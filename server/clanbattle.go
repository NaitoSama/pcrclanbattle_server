package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
	"pcrclanbattle_server/model"
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

func renewBoss(renewBoss db.Boss) error {
	// renew database bosses
	result := db.DB.Model(db.Boss{}).Where("id = ?", renewBoss.ID).Updates(renewBoss)
	if result.Error != nil {
		return result.Error
	}

	// renew cache
	lock.Lock()
	defer lock.Unlock()
	db.Cache.Bosses[renewBoss.ID-1] = renewBoss
	// broadcast
	broadcastData, _ := json.Marshal(renewBoss)
	Server.broadcast <- broadcastData
	return nil
}

func AttackBoss(message []byte, name string) error {
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
	lock.RLock()
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
	lock.RUnlock()
	if bossNewValue == 0 && bossNewStage == 0 && bossNewRound == 0 {
		return errors.New("invalid boss_id")
	}

	// renew database bosses and records
	newBoss := beforeAttackBoss
	newBoss.Round = bossNewRound
	newBoss.Stage = bossNewStage
	newBoss.Value = bossNewValue
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
		AttackFrom: name,
		AttackTo:   fmt.Sprintf("%d,%d,%d,%d", data.BossID, beforeAttackBoss.Stage, beforeAttackBoss.Round, beforeAttackBoss.Value),
		Damage:     common.If(data.AType == 1, beforeAttackBoss.Value, data.Value).(int64),
	}
	result := db.DB.Model(db.Record{}).Create(&record)
	if result.Error != nil {
		return result.Error
	}
	// renew cache
	lock.Lock()
	defer lock.Unlock()
	db.Cache.Records = append(db.Cache.Records, record)
	return errors.New("ok")
}
