package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/db"
)

func UploadUserPic(c *gin.Context) {
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
	usernameA, _ := c.Get("username")
	username, _ := usernameA.(string)

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "can not read file",
		})
		return
	}
	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "can not read file",
		})
		return
	}
	ETag, _ := common.CalculateETagForBytes(fileData)

	dst := fmt.Sprintf("./pic/%s.jpg", ETag)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "only jpg is supported",
		})
		return
	}

	dst16 := fmt.Sprintf("./pic/%s_128.jpg", ETag)
	err = pic16Gen(dst, dst16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "failed to compress original pic",
		})
		return
	}

	cacheUser := *db.Cache.Users[username]
	cacheUser.UserPic = dst
	cacheUser.PicETag = ETag
	cacheUser.UserPic16 = dst16
	cacheUser.Pic16ETag = ETag + "_128"
	err = renewUser(&cacheUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "failed to renew user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}

// pic16Gen compress the original image to 128 * 128
func pic16Gen(picPath string, dst string) error {
	file, err := os.Open(picPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// decode pic
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	// zoom image to 16x16 size
	newImg := resize.Resize(128, 128, img, resize.Lanczos3)

	// create output image file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// save images in JPEG format to output file
	jpeg.Encode(out, newImg, nil)
	return nil
}

func renewUser(newUser *db.User) error {
	result := db.DB.Model(db.User{}).Where("name = ?", newUser.Name).Updates(newUser)
	if result.Error != nil {
		return result.Error
	}
	lock.Lock()
	defer lock.Unlock()
	db.Cache.Users[newUser.Name] = newUser
	return nil
}
