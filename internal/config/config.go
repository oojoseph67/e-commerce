package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
	Admin    AdminConfig
}

type AdminConfig struct {
	DomainName string
	AdminCode  string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret              string
	ExpiresIn           time.Duration
	RefreshTokenExpires time.Duration
}

type AWSConfig struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	S3Bucket        string
	S3Endpoint      string
	UploadStorage   string
}

type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := getEnv("PORT", "1234")
	ginMode := getEnv("GIN_MODE", "development")

	dbHost := getEnv("DB_HOST", "host")
	dbPort := getEnv("DB_PORT", "1234")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "ecommerce")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	jwtSecret := getEnv("JWT_SECRET", "your_jwt_secret")
	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	refreshTokenExpires, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRES_IN", "72h"))

	awsRegion := getEnv("AWS_REGION", "your_jwt_secret")
	awsAccessKeyId := getEnv("AWS_ACCESS_KEY_ID", "test")
	awsSecretAccessKeyId := getEnv("AWS_SECRET_ACCESS_KEY_ID", "test")
	awsS3Bucket := getEnv("AWS_S3_BUCKET", "uploads")
	awsS3Endpoint := getEnv("AWS_S3_ENDPOINT", "http://localhost:9999")

	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	uploadStorage := getEnv("UPLOAD_STORAGE", "local")

	domainName := getEnv("COMPANY_DOMAIN", "example.com")
	adminCode := getEnv("ADMIN_CODE", "code")

	return &Config{
		Server: ServerConfig{
			Port:    port,
			GinMode: ginMode,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
		JWT: JWTConfig{
			ExpiresIn:           jwtExpiresIn,
			Secret:              jwtSecret,
			RefreshTokenExpires: refreshTokenExpires,
		},
		AWS: AWSConfig{
			Region:          awsRegion,
			AccessKeyId:     awsAccessKeyId,
			SecretAccessKey: awsSecretAccessKeyId,
			S3Bucket:        awsS3Bucket,
			S3Endpoint:      awsS3Endpoint,
			UploadStorage:   uploadStorage,
		},
		Upload: UploadConfig{
			Path:        uploadPath,
			MaxFileSize: maxUploadSize,
		},
		Admin: AdminConfig{
			DomainName: domainName,
			AdminCode:  adminCode,
		},
	}, nil

}
