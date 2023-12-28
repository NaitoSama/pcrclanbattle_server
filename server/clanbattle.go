package server

import (
	"encoding/json"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
	"pcrclanbattle_server/model"
)

// bossDefaultValue get value of a certain boss in a certain round
func bossDefaultValue(bossID int, round int) int64 {
	var bossConf = config.Config.Boss
	if round >= bossConf.StageSwitchRound[4] {
		return bossConf.StageSix[bossID-1]
	} else if round >= bossConf.StageSwitchRound[3] {
		return bossConf.StageFive[bossID-1]
	} else if round >= bossConf.StageSwitchRound[2] {
		return bossConf.StageFour[bossID-1]
	} else if round >= bossConf.StageSwitchRound[1] {
		return bossConf.StageThree[bossID-1]
	} else if round >= bossConf.StageSwitchRound[0] {
		return bossConf.StageTwo[bossID-1]
	} else if round < bossConf.StageSwitchRound[0] && round > 0 {
		return bossConf.StageOne[bossID-1]
	}
	return -1
}

func AttackBoss(message []byte) error {
	var bossNewRound int
	var bossNewValue int64
	data := model.AttackPayload{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		return err
	}
	lock.Lock()
	for i := 0; i < len(db.Cache.Bosses); i++ {
		if db.Cache.Bosses[i].ID == data.BossID {
			// defeat or damage boss
			if data.Value >= db.Cache.Bosses[i].Value {
				data.Type = 1
			} else {
				data.Type = 0
			}
			// defeat boss
			if data.Type == 1 {
				bossNewRound = db.Cache.Bosses[i].Round + 1
				bossNewValue = bossDefaultValue(db.Cache.Bosses[i].ID, bossNewRound)
			} else {
				// damage boss
				bossNewRound = db.Cache.Bosses[i].Round
				bossNewValue = db.Cache.Bosses[i].Value - data.Value
			}
			db.Cache.Bosses[i].Round = bossNewRound
			db.Cache.Bosses[i].Value = bossNewValue
			break
		}
	}
	lock.Unlock()

	// renew database bosses and records
	// broadcast
	return nil
}
