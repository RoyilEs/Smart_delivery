package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPassword 加密密码
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// ComparePasswords 验证密码 hash之后的密码  输入的密码
func ComparePasswords(hashPassword string, plainPassword string) bool {
	byteHash := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))
	if err != nil {
		return false
	}
	return true
}
