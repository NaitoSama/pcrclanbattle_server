package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("pekoStation")

// RandTokenSet this middleware will give a cookie to client
func RandTokenSet(c *gin.Context) {
	_, err := c.Cookie("pekoToken")
	if err != nil {
		c.SetCookie(
			"pekoToken",
			RandStringBytes(16),
			2678400,
			"/",
			"",
			false,
			true,
		)
	}
}

// JWTAuthentication this middleware will verify user token
func JWTAuthentication(c *gin.Context) {
	token, err := c.Cookie("pekoToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": "request with no token",
		})
		return
	}
	userID, username, userAuthority, ok := ParseJWT(token)
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"result": "you do not have permission to access this page",
		})
		return
	}
	userIDStr := fmt.Sprint(userID)
	userAuthorityStr := fmt.Sprint(userAuthority)
	c.Set("username", username)
	c.Set("user_id", userIDStr)
	c.Set("user_authority", userAuthorityStr)
}

type MyClaims struct {
	UserID        int    `json:"user_id"`
	UserName      string `json:"username"`
	UserAuthority int    `json:"user_authority"`
	jwt.StandardClaims
}

// NewJWT generate a new token
func NewJWT(userID int, username string, userAuthority int) (string, error) {
	claims := MyClaims{
		UserID:        userID,
		UserName:      username,
		UserAuthority: userAuthority,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,
			ExpiresAt: time.Now().Unix() + +60*60*24*31,
			Issuer:    "peko",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("Token Generates Failed!")
		return "", err
	}
	return tokenString, nil
}

// ParseJWT parse a token to userinfo
func ParseJWT(token string) (int, string, int, bool) {
	parseToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Println("ParseToken Failed!")
		return -1, "", 0, false
	}
	username := parseToken.Claims.(*MyClaims).UserName
	userID := parseToken.Claims.(*MyClaims).UserID
	userAuthority := parseToken.Claims.(*MyClaims).UserAuthority
	return userID, username, userAuthority, true
}

// OKTokenSet set user cookie
func OKTokenSet(c *gin.Context, token string) {
	c.SetCookie(
		"pekoToken",
		token,
		2678400,
		"/",
		"",
		false,
		true,
	)
}
