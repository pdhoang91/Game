# Oden - API Documentation

## Base URL

All API endpoints are available at: `https://api.oden-game.com/v1`

## Authentication

Most endpoints require authentication. Use the login endpoint to obtain a JWT token, then include it in all subsequent requests in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Authentication

#### Register New User

```
POST /auth/register
```

Request body:
```json
{
  "username": "player123",
  "email": "player@example.com",
  "password": "securePassword123"
}
```

Response:
```json
{
  "success": true,
  "user_id": "user_123456",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login

```
POST /auth/login
```

Request body:
```json
{
  "username": "player123",
  "password": "securePassword123"
}
```

Response:
```json
{
  "success": true,
  "user_id": "user_123456",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Heroes

#### Get Hero Collection

```
GET /heroes/list
```

Response:
```json
{
  "heroes": [
    {
      "id": "hero_12345",
      "hero_type_id": "warrior_01",
      "name": "Warrior",
      "level": 5,
      "experience": 220,
      "hp": 500,
      "atk": 50,
      "skills": [
        {
          "id": "skill_001",
          "name": "Mighty Slash",
          "description": "Deal 150% ATK to a single enemy",
          "damage_multiplier": 1.5,
          "cooldown": 3
        }
      ]
    },
    {
      "id": "hero_12346",
      "hero_type_id": "mage_01",
      "name": "Mage",
      "level": 3,
      "experience": 120,
      "hp": 300,
      "atk": 70,
      "skills": [
        {
          "id": "skill_002",
          "name": "Fireball",
          "description": "Deal 120% ATK to all enemies",
          "damage_multiplier": 1.2,
          "cooldown": 4
        }
      ]
    }
  ]
}
```

#### Summon Heroes

```
POST /heroes/summon
```

Request body:
```json
{
  "summon_type": "basic"  // or "premium"
}
```

Response:
```json
{
  "success": true,
  "heroes": [
    {
      "id": "hero_12347",
      "hero_type_id": "archer_01",
      "name": "Archer",
      "level": 1,
      "experience": 0,
      "hp": 250,
      "atk": 60,
      "skills": [
        {
          "id": "skill_003",
          "name": "Quick Shot",
          "description": "Deal 130% ATK to a single enemy",
          "damage_multiplier": 1.3,
          "cooldown": 2
        }
      ]
    }
  ]
}
```

### Team Management

#### Save Team Formation

```
POST /team/save
```

Request body:
```json
{
  "positions": {
    "1": "hero_12345",
    "2": "hero_12346",
    "3": "hero_12347",
    "4": null,
    "5": null
  }
}
```

Response:
```json
{
  "success": true,
  "team_id": "team_7890"
}
```

#### Get Team Formation

```
GET /team/get
```

Response:
```json
{
  "team_id": "team_7890",
  "positions": {
    "1": {
      "hero_id": "hero_12345",
      "hero_type_id": "warrior_01",
      "name": "Warrior",
      "level": 5
    },
    "2": {
      "hero_id": "hero_12346",
      "hero_type_id": "mage_01",
      "name": "Mage",
      "level": 3
    },
    "3": {
      "hero_id": "hero_12347",
      "hero_type_id": "archer_01",
      "name": "Archer",
      "level": 1
    },
    "4": null,
    "5": null
  }
}
```

### Battle System

#### Start Battle

```
POST /battle/start
```

Request body:
```json
{
  "stage_id": "stage_001"
}
```

Response:
```json
{
  "battle_id": "battle_4567",
  "result": "victory",
  "battle_log": [
    {
      "turn": 1,
      "actions": [
        {
          "actor": "hero_12345",
          "target": "enemy_001",
          "skill_used": "skill_001",
          "damage_dealt": 75,
          "target_hp_remaining": 125
        },
        {
          "actor": "enemy_001",
          "target": "hero_12345",
          "skill_used": "basic_attack",
          "damage_dealt": 30,
          "target_hp_remaining": 470
        }
      ]
    },
    // More turns...
  ],
  "rewards": {
    "gold": 100,
    "experience": {
      "hero_12345": 50,
      "hero_12346": 50,
      "hero_12347": 50
    },
    "items": []
  }
}
```

### Idle Rewards

#### Get Idle Rewards

```
GET /idle/rewards
```

Response:
```json
{
  "time_away": 7200,  // seconds
  "rewards": {
    "gold": 240,
    "experience": {
      "hero_12345": 30,
      "hero_12346": 30,
      "hero_12347": 30
    }
  }
}
```

#### Claim Idle Rewards

```
POST /idle/claim
```

Response:
```json
{
  "success": true,
  "time_claimed": "2025-03-21T00:01:30Z",
  "rewards": {
    "gold": 240,
    "experience": {
      "hero_12345": 30,
      "hero_12346": 30,
      "hero_12347": 30
    }
  }
}
```

## Error Responses

All endpoints return error responses in the following format:

```json
{
  "success": false,
  "error": "error_code",
  "message": "Human-readable error message"
}
```

Common error codes:
- `auth_required`: Authentication required
- `invalid_credentials`: Invalid username or password
- `resource_not_found`: Requested resource not found
- `insufficient_resources`: Not enough resources to perform action
- `server_error`: Internal server error 