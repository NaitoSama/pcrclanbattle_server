package common

import (
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

// PasswordEncryption Use hash and base64 with salt
func PasswordEncryption(password string) string {
	passwordB := []byte(password)
	salt := []byte("peko")
	hash := argon2.IDKey([]byte(passwordB), salt, 1, 64*1024, 4, 32)
	hashedPassword := base64.RawStdEncoding.EncodeToString(hash)
	return hashedPassword
}
