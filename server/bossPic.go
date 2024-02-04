package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

func UploadBossPic(c *gin.Context) {
	authority, _ := c.Get("user_authority")
	if authority.(int) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"result": "insufficient permissions",
		})
		return
	}
	file, err := c.FormFile("pic")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "need pic",
		})
		return
	}
	if filepath.Ext(file.Filename) != ".jpg" {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "only jpg is supported",
		})
		return
	}
	bossID, err := strconv.ParseInt(c.Param("boss"), 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "param boss is not int",
		})
		return
	}
	dst := fmt.Sprintf("./pic/%d.jpg", bossID)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "only jpg is supported",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
