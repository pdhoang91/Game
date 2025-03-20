# Oden - Server Setup Guide

This guide will walk you through setting up the Golang backend server for the Oden Idle RPG game.

## Prerequisites

- Go 1.18 or newer
- MySQL 8.0 or compatible database
- Git
- Basic knowledge of Go and RESTful APIs
- AWS account (for deployment)

## Development Environment Setup

### 1. Install Go

#### For Windows:
1. Download the Windows installer from [golang.org/dl/](https://golang.org/dl/)
2. Run the installer and follow the prompts
3. Verify installation by opening a command prompt and running:
   ```bash
   go version
   ```

#### For macOS:
1. Install using Homebrew:
   ```bash
   brew install go
   ```
2. Verify installation:
   ```bash
   go version
   ```

#### For Linux:
1. Download the tarball:
   ```bash
   wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
   ```
2. Extract it:
   ```bash
   sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
   ```
3. Add Go to your PATH in `~/.profile`:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   ```
4. Verify installation:
   ```bash
   go version
   ```

### 2. Install MySQL

#### For Windows:
1. Download MySQL installer from [mysql.com](https://dev.mysql.com/downloads/installer/)
2. Run the installer and follow the prompts
3. Make sure to note your root password

#### For macOS:
1. Install using Homebrew:
   ```bash
   brew install mysql
   ```
2. Start MySQL service:
   ```bash
   brew services start mysql
   ```

#### For Linux:
1. Install MySQL:
   ```bash
   sudo apt-get update
   sudo apt-get install mysql-server
   ```
2. Secure the installation:
   ```bash
   sudo mysql_secure_installation
   ```

### 3. Create Database

1. Log in to MySQL:
   ```bash
   mysql -u root -p
   ```

2. Create the database:
   ```sql
   CREATE DATABASE oden;
   ```

3. Create a user for the application:
   ```sql
   CREATE USER 'oden_user'@'localhost' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON oden.* TO 'oden_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

## Project Setup

### 1. Clone the Repository

```bash
git clone https://your-repository-url/oden.git
cd oden/server
```

### 2. Initialize Go Module

If not already initialized:

```bash
go mod init github.com/yourusername/oden
```

### 3. Install Dependencies

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql
go get -u github.com/golang-jwt/jwt/v4
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto/bcrypt
```

### 4. Configure the Application

Create a configuration file:

```bash
cp internal/config/config.example.json internal/config/config.json
```

Edit `config.json` with your settings:

```json
{
  "server": {
    "port": 8080,
    "host": "localhost",
    "environment": "development"
  },
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "oden_user",
    "password": "your_password",
    "name": "oden",
    "max_open_conns": 10,
    "max_idle_conns": 5
  },
  "auth": {
    "jwt_secret": "your-secret-key",
    "token_expiry": 72
  },
  "game": {
    "idle_reward_rate": 0.5,
    "max_team_size": 5
  }
}
```

## Project Structure

The server code is organized as follows:

```
server/
├── cmd/
│   ├── api/           # Main API application entry point
│   │   └── main.go
│   └── migrations/    # Database migration tool
│       └── main.go
├── internal/
│   ├── api/           # API handlers
│   │   ├── auth.go    # Authentication handlers
│   │   ├── battle.go  # Battle handlers
│   │   ├── hero.go    # Hero handlers
│   │   ├── idle.go    # Idle rewards handlers
│   │   ├── router.go  # API router setup
│   │   └── team.go    # Team handlers
│   ├── auth/          # Authentication logic
│   │   ├── jwt.go     # JWT token handling
│   │   └── password.go # Password hashing
│   ├── config/        # Configuration
│   │   ├── config.go
│   │   └── config.example.json
│   ├── db/            # Database interactions
│   │   ├── db.go      # Database connection
│   │   └── migrations/ # SQL migration scripts
│   │       └── 001_initial_schema.sql
│   ├── game/          # Game logic
│   │   ├── battle.go  # Battle calculations
│   │   ├── hero.go    # Hero management
│   │   ├── idle.go    # Idle rewards calculation
│   │   └── team.go    # Team management
│   └── model/         # Data models
│       ├── battle.go
│       ├── hero.go
│       ├── team.go
│       └── user.go
└── pkg/               # Public packages
    └── util/          # Utility functions
        └── random.go  # Random number generation
```

## Database Schema Setup

Create a database migration script at `internal/db/migrations/001_initial_schema.sql`:

```sql
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
);

-- Hero types table (template for heroes)
CREATE TABLE IF NOT EXISTS hero_types (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    rarity ENUM('common', 'rare', 'epic', 'legendary') NOT NULL,
    base_hp INT NOT NULL,
    base_atk INT NOT NULL,
    description TEXT,
    image_url VARCHAR(255)
);

-- Skills table
CREATE TABLE IF NOT EXISTS skills (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    damage_multiplier FLOAT NOT NULL,
    cooldown INT NOT NULL,
    targets_all BOOLEAN NOT NULL DEFAULT FALSE
);

-- Hero type skills mapping
CREATE TABLE IF NOT EXISTS hero_type_skills (
    hero_type_id VARCHAR(36) NOT NULL,
    skill_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (hero_type_id, skill_id),
    FOREIGN KEY (hero_type_id) REFERENCES hero_types(id),
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

-- Player's heroes
CREATE TABLE IF NOT EXISTS heroes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    hero_type_id VARCHAR(36) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    experience INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (hero_type_id) REFERENCES hero_types(id),
    INDEX idx_user_id (user_id)
);

-- Teams
CREATE TABLE IF NOT EXISTS teams (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    position_1 VARCHAR(36),
    position_2 VARCHAR(36),
    position_3 VARCHAR(36),
    position_4 VARCHAR(36),
    position_5 VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (position_1) REFERENCES heroes(id),
    FOREIGN KEY (position_2) REFERENCES heroes(id),
    FOREIGN KEY (position_3) REFERENCES heroes(id),
    FOREIGN KEY (position_4) REFERENCES heroes(id),
    FOREIGN KEY (position_5) REFERENCES heroes(id),
    INDEX idx_user_id (user_id)
);

-- Enemy types
CREATE TABLE IF NOT EXISTS enemy_types (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    hp INT NOT NULL,
    atk INT NOT NULL,
    description TEXT,
    image_url VARCHAR(255)
);

-- Stages
CREATE TABLE IF NOT EXISTS stages (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    enemy_1 VARCHAR(36),
    enemy_2 VARCHAR(36),
    enemy_3 VARCHAR(36),
    enemy_4 VARCHAR(36),
    enemy_5 VARCHAR(36),
    gold_reward INT NOT NULL,
    exp_reward INT NOT NULL,
    FOREIGN KEY (enemy_1) REFERENCES enemy_types(id),
    FOREIGN KEY (enemy_2) REFERENCES enemy_types(id),
    FOREIGN KEY (enemy_3) REFERENCES enemy_types(id),
    FOREIGN KEY (enemy_4) REFERENCES enemy_types(id),
    FOREIGN KEY (enemy_5) REFERENCES enemy_types(id)
);

-- Battle results
CREATE TABLE IF NOT EXISTS battle_results (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    team_id VARCHAR(36) NOT NULL,
    stage_id VARCHAR(36) NOT NULL,
    result ENUM('victory', 'defeat') NOT NULL,
    rewards JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (stage_id) REFERENCES stages(id),
    INDEX idx_user_id (user_id)
);

-- Player resources
CREATE TABLE IF NOT EXISTS player_resources (
    user_id VARCHAR(36) PRIMARY KEY,
    gold INT NOT NULL DEFAULT 0,
    premium_currency INT NOT NULL DEFAULT 0,
    last_idle_claim TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Insert initial hero types
INSERT INTO hero_types (id, name, rarity, base_hp, base_atk, description) VALUES
('warrior_01', 'Warrior', 'common', 500, 50, 'A strong warrior with high HP'),
('mage_01', 'Mage', 'common', 300, 70, 'A powerful mage with high ATK'),
('archer_01', 'Archer', 'common', 350, 60, 'A swift archer with balanced stats'),
('tank_01', 'Tank', 'rare', 700, 30, 'A defensive tank with very high HP'),
('assassin_01', 'Assassin', 'rare', 250, 80, 'A deadly assassin with very high ATK');

-- Insert initial skills
INSERT INTO skills (id, name, description, damage_multiplier, cooldown, targets_all) VALUES
('skill_001', 'Mighty Slash', 'Deal 150% ATK to a single enemy', 1.5, 3, FALSE),
('skill_002', 'Fireball', 'Deal 120% ATK to all enemies', 1.2, 4, TRUE),
('skill_003', 'Quick Shot', 'Deal 130% ATK to a single enemy', 1.3, 2, FALSE),
('skill_004', 'Shield Bash', 'Deal 100% ATK to a single enemy and reduce their ATK', 1.0, 3, FALSE),
('skill_005', 'Backstab', 'Deal 180% ATK to a single enemy', 1.8, 4, FALSE);

-- Map skills to hero types
INSERT INTO hero_type_skills (hero_type_id, skill_id) VALUES
('warrior_01', 'skill_001'),
('mage_01', 'skill_002'),
('archer_01', 'skill_003'),
('tank_01', 'skill_004'),
('assassin_01', 'skill_005');

-- Insert initial enemy types
INSERT INTO enemy_types (id, name, hp, atk, description) VALUES
('goblin_01', 'Goblin', 200, 30, 'A weak goblin'),
('orc_01', 'Orc', 400, 40, 'A stronger orc'),
('troll_01', 'Troll', 600, 50, 'A tough troll');

-- Insert initial stages
INSERT INTO stages (id, name, description, enemy_1, enemy_2, enemy_3, gold_reward, exp_reward) VALUES
('stage_001', 'Forest Path', 'A peaceful forest with some goblins', 'goblin_01', 'goblin_01', NULL, 100, 50),
('stage_002', 'Dark Cave', 'A dark cave with orcs', 'goblin_01', 'orc_01', 'goblin_01', 150, 75),
('stage_003', 'Mountain Pass', 'A dangerous mountain pass with trolls', 'orc_01', 'troll_01', 'orc_01', 200, 100);
```

Run the migration to set up the database schema:

```bash
cd server
go run cmd/migrations/main.go
```

## Implementing Core Components

### 1. Database Connection

Create a file at `internal/db/db.go`:

```go
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/oden/internal/config"
)

var DB *sql.DB

// Initialize sets up the database connection
func Initialize(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	DB.SetConnMaxLifetime(time.Hour)

	// Check the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
```

### 2. Config Model

Create a file at `internal/config/config.go`:

```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port        int    `json:"port"`
		Host        string `json:"host"`
		Environment string `json:"environment"`
	} `json:"server"`

	Database struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		User         string `json:"user"`
		Password     string `json:"password"`
		Name         string `json:"name"`
		MaxOpenConns int    `json:"max_open_conns"`
		MaxIdleConns int    `json:"max_idle_conns"`
	} `json:"database"`

	Auth struct {
		JWTSecret   string `json:"jwt_secret"`
		TokenExpiry int    `json:"token_expiry"` // hours
	} `json:"auth"`

	Game struct {
		IdleRewardRate float64 `json:"idle_reward_rate"`
		MaxTeamSize    int     `json:"max_team_size"`
	} `json:"game"`
}

// Load loads configuration from a file
func Load(filePath string) (*Config, error) {
	// Default values
	cfg := &Config{}
	cfg.Server.Port = 8080
	cfg.Server.Host = "localhost"
	cfg.Server.Environment = "development"

	cfg.Database.MaxOpenConns = 10
	cfg.Database.MaxIdleConns = 5

	cfg.Auth.TokenExpiry = 72

	cfg.Game.IdleRewardRate = 0.5
	cfg.Game.MaxTeamSize = 5

	// Load from file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
}
```

### 3. Main Entry Point

Create a file at `cmd/api/main.go`:

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/oden/internal/api"
	"github.com/yourusername/oden/internal/config"
	"github.com/yourusername/oden/internal/db"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "internal/config/config.json", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	err = db.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set up API server
	router := api.SetupRouter(cfg)

	// Start server in a goroutine
	go func() {
		serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server at %s...", serverAddr)

		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
```

### 4. API Router

Create a file at `internal/api/router.go`:

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/oden/internal/config"
)

// SetupRouter initializes the API routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())

	// API version grouping
	v1 := router.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", registerHandler)
			auth.POST("/login", loginHandler)
		}

		// Protected routes
		api := v1.Group("/")
		api.Use(authMiddleware(cfg))
		{
			// Heroes
			heroes := api.Group("/heroes")
			{
				heroes.GET("/list", getHeroesHandler)
				heroes.POST("/summon", summonHeroHandler)
			}

			// Team
			team := api.Group("/team")
			{
				team.GET("/get", getTeamHandler)
				team.POST("/save", saveTeamHandler)
			}

			// Battle
			battle := api.Group("/battle")
			{
				battle.POST("/start", startBattleHandler)
			}

			// Idle
			idle := api.Group("/idle")
			{
				idle.GET("/rewards", getIdleRewardsHandler)
				idle.POST("/claim", claimIdleRewardsHandler)
			}
		}
	}

	return router
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
```

## Running the Server

To run the server locally for development:

```bash
cd server
go run cmd/api/main.go
```

By default, the server will run on `http://localhost:8080`.

## Next Steps

After setting up the core server structure:

1. Implement the API handlers for each endpoint
2. Implement the authentication middleware
3. Create the game logic for battles, hero summoning, and idle rewards
4. Add unit tests for critical functionality
5. Set up continuous integration for automated testing

For deployment to AWS, refer to the [Deployment Guide](deployment.md).

## Troubleshooting

- **Database Connection Issues**: Verify your MySQL credentials and that the server is running
- **Missing Dependencies**: Run `go mod tidy` to install/update dependencies
- **Permission Issues**: Ensure your database user has sufficient privileges
- **Go Build Errors**: Make sure your Go version is 1.18 or newer 