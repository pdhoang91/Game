# Oden - System Architecture

## Overview

Oden is an Idle RPG mobile game with a client-server architecture. The client is built with Unity for iOS, while the backend is written in Golang and deployed on AWS.

## Architecture Diagram

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│             │         │             │         │             │
│  iOS Client │ ◄─────► │  API Server │ ◄─────► │  Database   │
│   (Unity)   │         │  (Golang)   │         │   (MySQL)   │
│             │         │             │         │             │
└─────────────┘         └─────────────┘         └─────────────┘
                              │
                              │
                              ▼
                        ┌─────────────┐
                        │             │
                        │ Asset Store │
                        │    (S3)     │
                        │             │
                        └─────────────┘
```

## Client Architecture (Unity)

### Scenes
- **LoginScene**: User authentication
- **MainScene**: Hub for hero management, team formation
- **BattleScene**: Auto-battle visualization

### Key Components
- **AuthManager**: Handles login/registration with backend
- **GameManager**: Central game state management
- **HeroManager**: Hero collection and team management
- **BattleManager**: Battle simulation and visualization
- **IdleManager**: Calculates and applies offline progress
- **APIClient**: Communication with backend API

## Server Architecture (Golang)

### Components
- **API Layer**: RESTful endpoints for client communication
- **Auth Service**: User authentication and session management
- **Game Logic**: Battle calculations, hero stats, team formation
- **Idle Processing**: Offline reward calculations
- **Database Layer**: Data persistence and retrieval

### API Endpoints
- `/auth/register`: New user registration
- `/auth/login`: User authentication
- `/heroes/list`: Get user's hero collection
- `/heroes/summon`: Summon new heroes
- `/team/save`: Save team formation
- `/battle/start`: Initiate battle
- `/battle/rewards`: Claim battle rewards
- `/idle/rewards`: Calculate and claim idle rewards

## Database Schema

### Users Table
- `id`: Unique user ID
- `username`: User's login name
- `password_hash`: Hashed user password
- `created_at`: Account creation timestamp
- `last_login`: Last login timestamp

### Heroes Table
- `id`: Unique hero ID
- `user_id`: Owner user ID
- `hero_type_id`: Reference to hero template
- `level`: Hero level
- `experience`: Current experience points
- `created_at`: Acquisition timestamp

### HeroTypes Table
- `id`: Unique hero type ID
- `name`: Hero name
- `rarity`: Hero rarity (common, rare, epic)
- `base_hp`: Base health points
- `base_atk`: Base attack power
- `skill_id`: Hero skill reference

### Teams Table
- `id`: Unique team ID
- `user_id`: Owner user ID
- `position_1`: Hero ID for position 1
- `position_2`: Hero ID for position 2
- `position_3`: Hero ID for position 3
- `position_4`: Hero ID for position 4
- `position_5`: Hero ID for position 5

### BattleResults Table
- `id`: Unique battle ID
- `user_id`: User ID
- `team_id`: Team ID used
- `enemy_team_id`: Enemy team ID
- `result`: Win/Loss
- `rewards`: JSON of earned rewards
- `timestamp`: Battle timestamp

## AWS Infrastructure

- **EC2**: Hosts the Golang API server
- **RDS MySQL**: Stores game data
- **S3**: Stores static assets like hero images
- **CloudFront**: CDN for asset delivery
- **CloudWatch**: Monitoring and alerts
- **Route 53**: DNS management 