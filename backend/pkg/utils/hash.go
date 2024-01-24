package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUuid() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
