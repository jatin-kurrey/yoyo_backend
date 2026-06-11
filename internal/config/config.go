package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                  string
	Host                    string
	Port                    string
	DatabaseURL             string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	DBSSLMode               string
	CORSAllowedOrigins      []string
	TrustedProxies          []string
	JWTSecret               string
	JWTAccessTokenTTL       time.Duration
	BcryptCost              int
	AutoMigrate             bool
	RequestBodyLimitBytes   int64
	RateLimitPerMinute      int
	LoginRateLimitPerHour   int
	ContactRateLimitPerHour int
	PaymentRateLimitPerHour int
	RazorpayKeyID           string
	RazorpayKeySecret       string
	RazorpayWebhookSecret   string
	RazorpayEnabled         bool
	AdminName               string
	AdminEmail              string
	AdminPassword           string
	UploadsStorage          string
	UploadDir               string
	MaxUploadSizeBytes      int64
	R2AccountID             string
	R2AccessKeyID           string
	R2SecretAccessKey       string
	R2Bucket                string
	R2PublicBaseURL         string
	R2Region                string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:                  getEnv("APP_ENV", "development"),
		Host:                    getEnv("HOST", "0.0.0.0"),
		Port:                    getEnv("PORT", "8080"),
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "postgres"),
		DBPassword:              getEnv("DB_PASSWORD", "postgres"),
		DBName:                  getEnv("DB_NAME", "yoyo_booking"),
		DBSSLMode:               getEnv("DB_SSLMODE", "disable"),
		CORSAllowedOrigins:      splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		TrustedProxies:          splitCSV(getEnv("TRUSTED_PROXIES", "")),
		JWTSecret:               getEnv("JWT_SECRET", "dev-change-this-secret"),
		JWTAccessTokenTTL:       getDuration("JWT_ACCESS_TOKEN_TTL", 24*time.Hour),
		BcryptCost:              getInt("BCRYPT_COST", 12),
		AutoMigrate:             getBool("AUTO_MIGRATE", true),
		RequestBodyLimitBytes:   int64(getInt("REQUEST_BODY_LIMIT_BYTES", 1<<20)),
		RateLimitPerMinute:      getInt("RATE_LIMIT_PER_MINUTE", 120),
		LoginRateLimitPerHour:   getInt("LOGIN_RATE_LIMIT_PER_HOUR", 10),
		ContactRateLimitPerHour: getInt("CONTACT_RATE_LIMIT_PER_HOUR", 20),
		PaymentRateLimitPerHour: getInt("PAYMENT_RATE_LIMIT_PER_HOUR", 30),
		RazorpayKeyID:           os.Getenv("RAZORPAY_KEY_ID"),
		RazorpayKeySecret:       os.Getenv("RAZORPAY_KEY_SECRET"),
		RazorpayWebhookSecret:   os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		RazorpayEnabled:         getBool("RAZORPAY_ENABLED", true),
		AdminName:               getEnv("ADMIN_NAME", "YOYO Super Admin"),
		AdminEmail:              os.Getenv("ADMIN_EMAIL"),
		AdminPassword:           os.Getenv("ADMIN_PASSWORD"),
		UploadsStorage:          getEnv("UPLOADS_STORAGE", "local"),
		UploadDir:               getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadSizeBytes:      int64(getInt("MAX_UPLOAD_SIZE_BYTES", 5<<20)), // Default 5MB
		R2AccountID:             os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:           os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey:       os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2Bucket:                os.Getenv("R2_BUCKET"),
		R2PublicBaseURL:         getEnv("R2_PUBLIC_BASE_URL", ""),
		R2Region:                getEnv("R2_REGION", "auto"),
	}

	if cfg.AppEnv == "production" {
		if cfg.JWTSecret == "" || cfg.JWTSecret == "dev-change-this-secret" {
			return nil, fmt.Errorf("JWT_SECRET must be set to a strong value in production")
		}
		if len(cfg.CORSAllowedOrigins) == 0 || cfg.CORSAllowedOrigins[0] == "*" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS cannot be wildcard in production")
		}
		if cfg.RazorpayEnabled && (cfg.RazorpayKeyID == "" || cfg.RazorpayKeySecret == "") {
			return nil, fmt.Errorf("Razorpay credentials are required when RAZORPAY_ENABLED=true")
		}
	}

	return cfg, nil
}

func (c Config) DatabaseDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}
