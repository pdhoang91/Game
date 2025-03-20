package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Auth     AuthConfig     `json:"auth"`
	Storage  StorageConfig  `json:"storage"`
	Game     GameConfig     `json:"game"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret   string `json:"jwt_secret"`
	TokenExpiry int    `json:"token_expiry"` // in hours
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UseSSL    bool   `json:"use_ssl"`
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	MaxIdleHours      int `json:"max_idle_hours"`
	IdleGoldPerMinute int `json:"idle_gold_per_minute"`
	IdleExpPerMinute  int `json:"idle_exp_per_minute"`
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	// Read configuration file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse configuration
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	overrideWithEnv(&config)

	return &config, nil
}

// overrideWithEnv overrides configuration with environment variables
func overrideWithEnv(config *Config) {
	// Database
	if host := os.Getenv("ODEN_DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("ODEN_DB_PORT"); port != "" {
		var portInt int
		if _, err := fmt.Sscanf(port, "%d", &portInt); err == nil {
			config.Database.Port = portInt
		}
	}
	if user := os.Getenv("ODEN_DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("ODEN_DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbName := os.Getenv("ODEN_DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}

	// Auth
	if jwtSecret := os.Getenv("ODEN_JWT_SECRET"); jwtSecret != "" {
		config.Auth.JWTSecret = jwtSecret
	}

	// Server
	if port := os.Getenv("ODEN_PORT"); port != "" {
		var portInt int
		if _, err := fmt.Sscanf(port, "%d", &portInt); err == nil {
			config.Server.Port = portInt
		}
	}

	// Storage
	if endpoint := os.Getenv("ODEN_STORAGE_ENDPOINT"); endpoint != "" {
		config.Storage.Endpoint = endpoint
	}
	if bucket := os.Getenv("ODEN_STORAGE_BUCKET"); bucket != "" {
		config.Storage.Bucket = bucket
	}
	if accessKey := os.Getenv("ODEN_STORAGE_ACCESS_KEY"); accessKey != "" {
		config.Storage.AccessKey = accessKey
	}
	if secretKey := os.Getenv("ODEN_STORAGE_SECRET_KEY"); secretKey != "" {
		config.Storage.SecretKey = secretKey
	}
} 