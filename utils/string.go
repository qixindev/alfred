package utils

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StructToString(obj any) string {
	res, err := make([]byte, 0), errors.New("")
	if res, err = json.Marshal(&obj); err != nil {
		return ""
	}
	return string(res)
}

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetDeviceUserCode() string {
	b := make([]byte, 9)
	for i := 0; i < 9; i++ {
		if i == 4 {
			b[i] = byte('-')
			i++
		}
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
