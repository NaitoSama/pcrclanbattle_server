package server

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
	"pcrclanbattle_server/model"
	"strings"
	"time"
)

// bossDefaultValue get stage and value of a certain boss in a certain round
func bossDefaultValue(bossID int, round int) (int, int64) {
	var bossConf = config.Config.Boss
	if round >= bossConf.StageSwitchRound[4] {
		return 6, bossConf.StageSix[bossID-1]
	} else if round >= bossConf.StageSwitchRound[3] {
		return 5, bossConf.StageFive[bossID-1]
	} else if round >= bossConf.StageSwitchRound[2] {
		return 4, bossConf.StageFour[bossID-1]
	} else if round >= bossConf.StageSwitchRound[1] {
		return 3, bossConf.StageThree[bossID-1]
	} else if round >= bossConf.StageSwitchRound[0] {
		return 2, bossConf.StageTwo[bossID-1]
	} else if round < bossConf.StageSwitchRound[0] && round > 0 {
		return 1, bossConf.StageOne[bossID-1]
	}
	return -1, -1
}

// renewBoss thread is not safe
func renewBoss(renewBoss db.Boss) error {
	// renew database bosses
	result := db.DB.Model(db.Boss{}).Where("id = ?", renewBoss.ID).Updates(renewBoss)
	if result.Error != nil {
		return result.Error
	}

	// renew cache
	db.Cache.Bosses[renewBoss.ID-1] = renewBoss
	// broadcast
	broadcastData, _ := json.Marshal(renewBoss)
	Server.broadcast <- broadcastData
	return nil
}

//// parseBossStatus will return a boss status with default
//func parseBossStatus(bossStatusStr string) db.Boss {
//	data := strings.Split(bossStatusStr, "|")
//	bossID, _ := strconv.Atoi(data[0])
//	stage, _ := strconv.Atoi(data[1])
//	round, _ := strconv.Atoi(data[2])
//	value, _ := strconv.ParseInt(data[3], 10, 64)
//	boss := db.Boss{
//		ID:      bossID,
//		Stage:   stage,
//		Round:   round,
//		Value:   value,
//		WhoIsIn: " ",
//	}
//	return boss
//}

// AttackBoss need attacker name
func AttackBoss(message []byte, name string) error {
	lock.Lock()
	defer lock.Unlock()
	var bossNewStage int
	var bossNewRound int
	var bossNewValue int64
	var beforeAttackBoss db.Boss
	data := model.AttackPayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}
	if data.Value < 0 {
		return errors.New("invalid value")
	}
	// boss info after attacking
	for i := 0; i < len(db.Cache.Bosses); i++ {
		if db.Cache.Bosses[i].ID == data.BossID {
			beforeAttackBoss = db.Cache.Bosses[i]
			// defeat or damage boss
			if data.Value >= db.Cache.Bosses[i].Value {
				data.AType = 1
			} else {
				data.AType = 0
			}

			if data.AType == 1 {
				// defeat boss
				bossNewRound = db.Cache.Bosses[i].Round + 1
				bossNewStage, bossNewValue = bossDefaultValue(db.Cache.Bosses[i].ID, bossNewRound)
			} else {
				// damage boss
				bossNewRound = db.Cache.Bosses[i].Round
				bossNewValue = db.Cache.Bosses[i].Value - data.Value
				bossNewStage = db.Cache.Bosses[i].Stage
			}
			break
		}
	}
	if bossNewValue == 0 && bossNewStage == 0 && bossNewRound == 0 {
		return errors.New("invalid boss_id")
	}

	// renew database bosses and records
	newBoss := beforeAttackBoss
	newBoss.Round = bossNewRound
	newBoss.Stage = bossNewStage
	newBoss.Value = bossNewValue
	if data.AType == 1 {
		newBoss.WhoIsIn = " "
		newBoss.Tree = " "
		newBoss.ValueD = bossNewValue
	} else if name == newBoss.WhoIsIn {
		newBoss.WhoIsIn = " "
	}
	err = renewBoss(newBoss)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	record := db.Record{
		Model: gorm.Model{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		AttackFrom:        name,
		AttackTo:          beforeAttackBoss.ID,
		Damage:            common.If(data.AType == 1, beforeAttackBoss.Value, data.Value).(int64),
		BeforeBossStage:   beforeAttackBoss.Stage,
		BeforeBossRound:   beforeAttackBoss.Round,
		BeforeBossValue:   beforeAttackBoss.Value,
		BeforeBossWhoIsIn: beforeAttackBoss.WhoIsIn,
		BeforeBossTree:    beforeAttackBoss.Tree,
		BeforeBossValueD:  beforeAttackBoss.ValueD,
		CanUndo:           1,
	}
	result := db.DB.Model(db.Record{}).Create(&record)
	if result.Error != nil {
		return result.Error
	}
	// renew cache
	db.Cache.Records = append(db.Cache.Records, record)
	broadcastData, _ := json.Marshal(record)
	Server.broadcast <- broadcastData
	return errors.New("ok")
}

// ReviseBoss by sending boss_id and round and value
func ReviseBoss(message []byte) error {
	lock.Lock()
	defer lock.Unlock()
	data := model.RevisePayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}
	bossNewStage, bossNewDefaultValue := bossDefaultValue(data.BossID, data.Round)
	newBoss := db.Boss{
		ID:      data.BossID,
		Stage:   bossNewStage,
		Round:   data.Round,
		Value:   common.If(data.Value > bossNewDefaultValue, bossNewDefaultValue, data.Value).(int64),
		WhoIsIn: " ",
		Tree:    " ",
		ValueD:  bossNewDefaultValue,
		PicETag: db.Cache.Bosses[data.BossID-1].PicETag,
	}
	err = renewBoss(newBoss)
	if err != nil {
		return err
	}
	return errors.New("ok")
}

func Undo(message []byte, name string) error {
	lock.Lock()
	defer lock.Unlock()
	var bossStatus db.Boss // it is the state of the boss we want to have after undo the attack
	var bossStatusDataIndex int
	var findBossStatus bool

	data := model.UndoPayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}
	recordsLen := len(db.Cache.Records)
	round := config.Config.ClanBattle.CanBeUndoRecordsUP

	for i := recordsLen - 1; i >= 0; i-- {
		round--
		if data.BossID == db.Cache.Records[i].AttackTo {
			if name == db.Cache.Records[i].AttackFrom && db.Cache.Records[i].CanUndo == 1 {
				db.Cache.Records[i].CanUndo = 0
				db.DB.Model(db.Record{}).Where("id = ?", db.Cache.Records[i].ID).Update("can_undo", 0)
				bossStatus = db.Boss{
					ID:      db.Cache.Records[i].AttackTo,
					Stage:   db.Cache.Records[i].BeforeBossStage,
					Round:   db.Cache.Records[i].BeforeBossRound,
					Value:   db.Cache.Records[i].BeforeBossValue,
					WhoIsIn: db.Cache.Records[i].BeforeBossWhoIsIn,
					Tree:    db.Cache.Records[i].BeforeBossTree,
					ValueD:  db.Cache.Records[i].BeforeBossValueD,
					PicETag: db.Cache.Bosses[db.Cache.Records[i].AttackTo-1].PicETag,
				}
				bossStatusDataIndex = i
				findBossStatus = true
				break
			} else {
				return errors.New("can not undo, because someone has attacked the boss")
			}
		}
		if round <= 0 {
			return errors.New("can not find " + name + "'s record in the latest records")
		}
	}
	if !findBossStatus {
		return errors.New(name + " has no record")
	}

	err = renewBoss(bossStatus)
	if err != nil {
		return err
	}
	// renew cache
	db.Cache.Records = append(db.Cache.Records[:bossStatusDataIndex], db.Cache.Records[bossStatusDataIndex+1:]...)
	return errors.New("ok")
}

func ImIn(message []byte, name string) error {
	lock.Lock()
	defer lock.Unlock()
	data := model.ImInPayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}

	newBoss := db.Cache.Bosses[data.BossID-1]
	if newBoss.WhoIsIn != name && newBoss.WhoIsIn != " " {
		return errors.New("someone is attacking the boss")
	}
	if newBoss.WhoIsIn == name {
		return errors.New("you are in boss now")
	}

	newBoss.WhoIsIn = name
	err = renewBoss(newBoss)
	if err != nil {
		return err
	}
	return errors.New("ok")
}

func ImOut(message []byte, name string) error {
	lock.Lock()
	defer lock.Unlock()
	data := model.ImOutPayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}

	newBoss := db.Cache.Bosses[data.BossID-1]
	if newBoss.WhoIsIn != name {
		return errors.New("you are not in boss now")
	}

	newBoss.WhoIsIn = " "
	err = renewBoss(newBoss)
	if err != nil {
		return err
	}
	return errors.New("ok")
}

func OnTree(message []byte, name string) error {
	lock.Lock()
	defer lock.Unlock()
	data := model.OnTreePayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}

	newBoss := db.Cache.Bosses[data.BossID-1]
	// if tree has somebody on it
	if newBoss.Tree != " " {
		treeArray := strings.Split(newBoss.Tree, "|")
		_, findResult := common.SliceFind(treeArray, name)
		if findResult {
			return errors.New("you are already on tree")
		}
		newBoss.Tree += "|" + name
	} else { // there is nobody on tree
		newBoss.Tree = name
	}
	err = renewBoss(newBoss)
	if err != nil {
		return err
	}
	return errors.New("ok")
}

func DownTree(message []byte, name string) error {
	lock.Lock()
	defer lock.Unlock()
	data := model.DownTreePayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}

	newBoss := db.Cache.Bosses[data.BossID-1]
	treeArray := strings.Split(newBoss.Tree, "|")
	index, isOnTree := common.SliceFind(treeArray, name)
	// if user is not on tree
	if !isOnTree {
		return errors.New("you are not on tree now")
	}
	newTreeArray := append(treeArray[:index], treeArray[index+1:]...)
	if len(newTreeArray) > 0 {
		newBoss.Tree = strings.Join(newTreeArray, "|")
	} else {
		newBoss.Tree = " "
	}

	err = renewBoss(newBoss)
	if err != nil {
		return err
	}
	return errors.New("ok")
}
