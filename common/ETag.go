package common

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
)

// 计算文件的 ETag
func CalculateETag(filename string) (string, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 计算文件内容的 SHA-1 哈希值
	hash := sha1.New()
	hash.Write(content)
	sha1Hash := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	etag := hex.EncodeToString(sha1Hash)

	return etag, nil
}
