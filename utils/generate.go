package utils

import (
	"math/rand"
	"time"
	"fmt"
	"crypto/md5"
)

// GenerateRandomString 根据传入的长度l，生成随机长度l的字符串
func GenerateRandomString(length uint) string  {
	var seeks = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+-,.;':|"
	sequences := make([]byte, 16)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range sequences {
		sequences[i] = seeks[r.Intn(len(seeks))]
	}
	return string(sequences)
}

// GenerateToken 生成Toekn
func GenerateToken() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(GenerateRandomString(16))))
}

func MD5(s string) string  {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// CryptPassword 进行MD%加密
func CryptPassword(origin string) string {
	return MD5(origin)
}