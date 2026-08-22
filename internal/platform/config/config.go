package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
	Admin    AdminConfig
	Log      LogConfig
}

type AppConfig struct {
	Name           string
	Env            string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestMaxSize int64
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type SecurityConfig struct {
	PasswordMinLength int
	LoginMaxAttempts  int
	RequireUpper      bool
	RequireDigit      bool
	RequireSpecial    bool
	LockoutDuration   time.Duration
	IdempotencyTTL    time.Duration
}

type AdminConfig struct {
	Username string
	Password string
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	config := &Config{
		App: AppConfig{Name: env("APP_NAME", "go-waste-routes"), Env: env("APP_ENV", "development"), Port: env("APP_PORT", "8080")},
		DB: DatabaseConfig{Host: env("DB_HOST", "localhost"), Port: env("DB_PORT", "3306"), User: env("DB_USER", "waste_user"), Password: env("DB_PASSWORD", "waste_pass_2026"), Name: env("DB_NAME", "go_waste_routes")},
		Redis: RedisConfig{Host: env("REDIS_HOST", "localhost"), Port: env("REDIS_PORT", "6379"), Password: env("REDIS_PASSWORD", "")},
		JWT: JWTConfig{Secret: env("JWT_SECRET", "your-256-bit-secret-change-in-production")},
		Admin: AdminConfig{Username: env("ADMIN_USERNAME", "admin"), Password: env("ADMIN_PASSWORD", "Admin123!")},
		Log: LogConfig{Level: env("LOG_LEVEL", "debug"), Format: env("LOG_FORMAT", "json")},
	}
	config.App.ReadTimeout = duration("APP_READ_TIMEOUT", "30s")
	config.App.WriteTimeout = duration("APP_WRITE_TIMEOUT", "30s")
	config.App.RequestMaxSize = int64Value("APP_REQUEST_MAX_SIZE", 10485760)
	config.DB.MaxOpenConns = intValue("DB_MAX_OPEN_CONNS", 25)
	config.DB.MaxIdleConns = intValue("DB_MAX_IDLE_CONNS", 10)
	config.DB.ConnMaxLifetime = duration("DB_CONN_MAX_LIFETIME", "300s")
	config.Redis.DB = intValue("REDIS_DB", 0)
	config.Redis.PoolSize = intValue("REDIS_POOL_SIZE", 20)
	config.JWT.AccessTTL = duration("JWT_ACCESS_TTL", "15m")
	config.JWT.RefreshTTL = duration("JWT_REFRESH_TTL", "168h")
	config.Security.PasswordMinLength = intValue("PASSWORD_MIN_LENGTH", 8)
	config.Security.LoginMaxAttempts = intValue("LOGIN_MAX_ATTEMPTS", 5)
	config.Security.RequireUpper = boolValue("PASSWORD_REQUIRE_UPPER", true)
	config.Security.RequireDigit = boolValue("PASSWORD_REQUIRE_DIGIT", true)
	config.Security.RequireSpecial = boolValue("PASSWORD_REQUIRE_SPECIAL", true)
	config.Security.LockoutDuration = duration("LOGIN_LOCKOUT_DURATION", "15m")
	config.Security.IdempotencyTTL = duration("IDEMPOTENCY_TTL", "24h")
	if len(config.JWT.Secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	return config, nil
}

func env(key, fallback string) string { if value := os.Getenv(key); value != "" { return value }; return fallback }
func duration(key, fallback string) time.Duration { value, err := time.ParseDuration(env(key, fallback)); if err != nil { value, _ = time.ParseDuration(fallback) }; return value }
func intValue(key string, fallback int) int { value, err := strconv.Atoi(env(key, strconv.Itoa(fallback))); if err != nil { return fallback }; return value }
func int64Value(key string, fallback int64) int64 { value, err := strconv.ParseInt(env(key, strconv.FormatInt(fallback, 10)), 10, 64); if err != nil { return fallback }; return value }
func boolValue(key string, fallback bool) bool { value, err := strconv.ParseBool(env(key, strconv.FormatBool(fallback))); if err != nil { return fallback }; return value }
