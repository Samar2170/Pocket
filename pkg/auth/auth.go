package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"pocket/pkg/auditlog"
	"pocket/pkg/db"
	"pocket/pkg/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type APIKey struct {
	*gorm.Model
	KeyHash string
}

var SecretKey string

func init() {
	envPath := utils.Basedir + "/.env"
	err := godotenv.Load(envPath)
	if err != nil {
		auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error loading .env file")
	}
	SecretKey = os.Getenv("SECRETKEY")
	db.DB.AutoMigrate(&APIKey{})
}

func GenerateKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		auditlog.AuditLogger.Error().Str("error", err.Error()).Msg("Error generating key")
		panic(err)
	}
	return hex.EncodeToString(key)
}

func GetNewKey() {
	key := GenerateKey()
	keyHash := HashKey(key)
	db.DB.Model(&APIKey{}).Where("id = ?", 1).Update("key_hash", keyHash)
	// db.DB.Create(&APIKey{KeyHash: keyHash})
	fmt.Println(key)
}

func HashKey(apiKey string) string {
	combined := append([]byte(apiKey), []byte(SecretKey)...)
	hash := sha256.New()
	hash.Write(combined)
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func IsKeyValid(key string) bool {
	keyHash := HashKey(key)
	var apiKey APIKey
	err := db.DB.Where("key_hash = ?", keyHash).First(&apiKey).Error
	return err == nil
}
