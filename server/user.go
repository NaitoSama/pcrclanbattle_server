package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
	"strings"
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
	if username == " " {
		c.JSON(http.StatusBadRequest, gin.H{"result": "username can not be blank"})
		return
	}
	if strings.Contains(username, "|") {
		c.JSON(http.StatusBadRequest, gin.H{"result": "username can not contains \"|\""})
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
	lock.Lock()
	defer lock.Unlock()
	db.Cache.Users[user.Name] = &user
	token, _ := common.NewJWT(0, user.Name, 0)
	common.OKTokenSet(c, token)
	c.JSON(http.StatusOK, gin.H{"result": "registered successfully"})
}

func GetUserInfoFromJWT(c *gin.Context) {
	data := make(map[string]string)
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	jwt, ok := data["jwt"]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	userID, username, userAuthority, ok := common.ParseJWT(jwt)
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":        userID,
		"username":       username,
		"user_authority": userAuthority,
	})
}

func ChangePassword(c *gin.Context) {
	json := make(map[string]string)
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "must be string"})
		return
	}
	username, usernameExists := json["username"]
	oldPassword, passwordExists := json["old_password"]
	newPassword, newPasswordExists := json["new_password"]
	if !usernameExists || !passwordExists || !newPasswordExists {
		c.JSON(http.StatusBadRequest, gin.H{"result": "Invalid JSON structure"})
		return
	}
	if oldPassword == newPassword {
		c.JSON(http.StatusBadRequest, gin.H{"result": "new password cannot be the same as the old one"})
		return
	}

	user, ok := db.Cache.Users[username]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"result": "user dose not exist"})
		return
	}
	newUser := *user
	if newUser.Password != common.PasswordEncryption(oldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"result": "wrong old password"})
		return
	}
	newUser.Password = common.PasswordEncryption(newPassword)
	err = renewUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": "can not update new password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "change password successfully"})
}
