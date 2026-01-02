package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConnections  int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// loads configuration from env variables
func LoadConfig() (*DBConfig, error) {
	// load .env file in development
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}
	maxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNECTIONS", "100"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNECTIONS", "10"))
	connMaxLifetime, _ := strconv.Atoi(getEnv("DB_CONNECTION_MAX_LIFETIME", "3600"))

	return &DBConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "militaryindex"),
		Password:        getEnv("DB_PASSWORD", "SecurePassword123"),
		DBName:          getEnv("DB_NAME", "military_index_db"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxConnections:  maxConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: time.Duration(connMaxLifetime) * time.Second,
	}, nil
}

// ConnectDatabase establishes database connection
func ConnectDatabase() (*gorm.DB, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Built DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	// Configure logger based on environment
	var logLevel logger.LogLevel
	if os.Getenv("APP_ENV") == "production" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	// Open Connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables prepared statements for Supabase pooler compatibility
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	log.Printf("Database connected successfully to %s %s %s",
		config.Host, config.Port, config.DBName)

	return db, nil
}

// HealthCheck checks database connection health
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
