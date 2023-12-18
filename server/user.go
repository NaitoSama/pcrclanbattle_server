package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
)

func Login(c *gin.Context) {
	json := make(map[string]string)
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "must be string"})
		return
	}
	username, usernameExists := json["username"]
	password, passwordExists := json["password"]

	if !usernameExists || !passwordExists {
		c.JSON(http.StatusBadRequest, gin.H{"result": "Invalid JSON structure"})
		return
	}
	user := db.User{}
	result := db.DB.Model(db.User{}).Where("name = ? and password = ?", username, common.PasswordEncryption(password)).First(&user)
	if result.RowsAffected == 0 || result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "wrong name or password"})
		return
	}
	token, _ := common.NewJWT(user.UserID, user.Name, user.Permission)
	common.OKTokenSet(c, token)
	c.JSON(http.StatusOK, gin.H{"result": user.Name + ",你好"})
}

func Register(c *gin.Context) {
	json := make(map[string]string)
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "must be string"})
		return
	}
	username, usernameExists := json["username"]
	password, passwordExists := json["password"]
	registerCode, registerCodeExists := json["register_code"]

	if !usernameExists || !passwordExists || !registerCodeExists {
		c.JSON(http.StatusBadRequest, gin.H{"result": "Invalid JSON structure"})
		return
	}
	if registerCode != config.Config.General.RegisterCode {
		c.JSON(http.StatusBadRequest, gin.H{"result": "invalid register code"})
		return
	}
	user := db.User{}
	result := db.DB.Model(user).Where("name = ?", username).First(&user)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"result": "username already existed"})
		return
	}

	user.Name = username
	user.Password = common.PasswordEncryption(password)
	result = db.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"result": "registration failed"})
		return
	}
	token, _ := common.NewJWT(0, user.Name, 0)
	common.OKTokenSet(c, token)
	c.JSON(http.StatusOK, gin.H{"result": "registered successfully"})
}
