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
	c.JSON(http.StatusOK, db.Cache.Records)
}

func GetBosses(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()
	c.JSON(http.StatusOK, db.Cache.Bosses)
}

func GetUsers(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()
	var users []db.User
	reqUsers := c.QueryArray("users")
	for i := 0; i < len(reqUsers); i++ {
		users = append(users, *db.Cache.Users[reqUsers[i]])
	}
	c.JSON(http.StatusOK, users)
}
