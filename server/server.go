package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pcrclanbattle_server/db"
)

func ServerInit() {}

func GetRecords(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()
	records := db.Cache.Records
	c.JSON(http.StatusOK, records)
}

func GetBosses(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()
	bosses := db.Cache.Bosses
	c.JSON(http.StatusOK, bosses)
}
