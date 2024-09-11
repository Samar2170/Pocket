package utils

import (
	"math/rand"
	"time"
)

func CheckArray[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}

var ToBytesSize = map[string]int{
	"bytes": 1,
	"kb":    1024,
	"mb":    1024 * 1024,
	"gb":    1024 * 1024 * 1024,
	"tb":    1024 * 1024 * 1024 * 1024,
}

func ConvertFileSize(size float64, currentUnit, resultUnit string) float64 {
	toByteMultiplier := ToBytesSize[currentUnit]
	result := size * float64(toByteMultiplier)
	resultUnitFactor := ToBytesSize[resultUnit]
	result = result / float64(resultUnitFactor)
	return result
}
