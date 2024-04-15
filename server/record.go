package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
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
