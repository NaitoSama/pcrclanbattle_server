package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
	"strconv"
	"time"
)

// RecordsArchiving will set records archived and reset bosses
func RecordsArchiving(c *gin.Context) {
	var records []db.Record
	authority, _ := c.Get("user_authority")
	if authority != "2" {
		c.JSON(http.StatusForbidden, gin.H{
			"result": "insufficient permissions",
		})
		return
	}
	archiveName := c.Query("archive_name")
	if archiveName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "must have a archive name",
		})
		return
	}
	db.DB.Find(&records, "archive_id = ?", "")
	db.DB.Model(&records).Updates(db.Record{ArchiveID: archiveName})
	lock.Lock()
	defer lock.Unlock()
	db.DB.Model(db.Record{}).Find(&db.Cache.Records)

	for i := 0; i < 5; i++ {
		boss := db.Boss{
			ID:      i + 1,
			Stage:   1,
			Round:   1,
			Value:   config.Config.Boss.StageOne[i],
			WhoIsIn: " ",
			Tree:    " ",
			ValueD:  config.Config.Boss.StageOne[i],
			PicETag: db.Cache.Bosses[i].PicETag,
		}
		_ = renewBoss(boss)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}

func DeleteRecords(c *gin.Context) {
	authority, _ := c.Get("user_authority")
	if authority == "0" {
		c.JSON(http.StatusForbidden, gin.H{
			"result": "insufficient permissions",
		})
		return
	}
	recordIDS := c.Query("record_id")
	recordID, err := strconv.ParseUint(recordIDS, 10, 64)
	if recordIDS == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "invalid record id",
		})
		return
	}

	db.DB.Model(db.Record{}).Where("id = ?", recordID).Update("deleted_at", time.Now())
	lock.Lock()
	defer lock.Unlock()
	length := len(db.Cache.Records)
	for i := 0; i < length; i++ {
		j := length - 1 - i
		if uint64(db.Cache.Records[j].ID) == recordID {
			toDeleteRecord := db.Cache.Records[j]
			var temp []db.Record
			if j != length-1 {
				temp = append(db.Cache.Records[:j], db.Cache.Records[j+1:]...)
			} else {
				temp = db.Cache.Records[:j]
			}
			db.Cache.Records = temp
			content := db.Content{Type: "record_delete", Data: toDeleteRecord}
			// broadcast
			broadcastData, _ := json.Marshal(content)
			Server.broadcast <- broadcastData
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
