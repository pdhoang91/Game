package model

import (
	"encoding/json"
	"time"
)

// BattleResult represents the result of a battle
type BattleResult struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TeamID    string    `json:"team_id"`
	StageID   string    `json:"stage_id"`
	Result    string    `json:"result"` // victory, defeat
	RewardsJSON string   `json:"-"`     // JSON string stored in DB
	Rewards   *Rewards  `json:"rewards,omitempty"` // Parsed from RewardsJSON
	CreatedAt time.Time `json:"created_at"`
	
	// Battle log for client-side visualization
	BattleLog []BattleTurn `json:"battle_log,omitempty"`
}

// BattleTurn represents a turn in a battle
type BattleTurn struct {
	Turn    int           `json:"turn"`
	Actions []BattleAction `json:"actions"`
}

// BattleAction represents an action in a battle turn
type BattleAction struct {
	Actor            string `json:"actor"`            // Hero or enemy ID
	Target           string `json:"target"`           // Hero or enemy ID
	SkillUsed        string `json:"skill_used"`       // Skill ID
	DamageDealt      int    `json:"damage_dealt"`     // Damage dealt
	TargetHPRemaining int   `json:"target_hp_remaining"` // Remaining HP of target
}

// Rewards represents rewards from a battle
type Rewards struct {
	Gold       int                `json:"gold"`
	Experience map[string]int     `json:"experience"` // Hero ID -> XP
	Items      []string           `json:"items"`      // Item IDs
}

// NewBattleResult creates a new battle result instance
func NewBattleResult(id, userID, teamID, stageID string, result string, rewards *Rewards) *BattleResult {
	return &BattleResult{
		ID:        id,
		UserID:    userID,
		TeamID:    teamID,
		StageID:   stageID,
		Result:    result,
		Rewards:   rewards,
		CreatedAt: time.Now(),
		BattleLog: make([]BattleTurn, 0),
	}
}

// Stage represents a stage in the game
type Stage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Enemy1      string `json:"enemy_1,omitempty"`
	Enemy2      string `json:"enemy_2,omitempty"`
	Enemy3      string `json:"enemy_3,omitempty"`
	Enemy4      string `json:"enemy_4,omitempty"`
	Enemy5      string `json:"enemy_5,omitempty"`
	GoldReward  int    `json:"gold_reward"`
	ExpReward   int    `json:"exp_reward"`
	
	// Computed fields (not stored in DB)
	Enemies    []*Enemy `json:"enemies,omitempty"`
}

// Enemy represents an enemy in a stage
type Enemy struct {
	ID          string `json:"id"`
	TypeID      string `json:"type_id"`
	Name        string `json:"name"`
	HP          int    `json:"hp"`
	ATK         int    `json:"atk"`
	Description string `json:"description,omitempty"`
	
	// Runtime battle state
	CurrentHP   int    `json:"-"`
}

// GetEnemyIDs returns all enemy IDs in the stage
func (s *Stage) GetEnemyIDs() []string {
	var ids []string
	if s.Enemy1 != "" {
		ids = append(ids, s.Enemy1)
	}
	if s.Enemy2 != "" {
		ids = append(ids, s.Enemy2)
	}
	if s.Enemy3 != "" {
		ids = append(ids, s.Enemy3)
	}
	if s.Enemy4 != "" {
		ids = append(ids, s.Enemy4)
	}
	if s.Enemy5 != "" {
		ids = append(ids, s.Enemy5)
	}
	return ids
}

// SetRewardsJSON converts the Rewards object to JSON and stores it
func (br *BattleResult) SetRewardsJSON() error {
	if br.Rewards == nil {
		br.RewardsJSON = "{}"
		return nil
	}
	
	jsonBytes, err := json.Marshal(br.Rewards)
	if err != nil {
		return err
	}
	
	br.RewardsJSON = string(jsonBytes)
	return nil
}

// ParseRewardsJSON parses the JSON rewards string
func (br *BattleResult) ParseRewardsJSON() error {
	if br.RewardsJSON == "" {
		br.Rewards = &Rewards{
			Gold:       0,
			Experience: make(map[string]int),
			Items:      []string{},
		}
		return nil
	}
	
	var rewards Rewards
	err := json.Unmarshal([]byte(br.RewardsJSON), &rewards)
	if err != nil {
		return err
	}
	
	br.Rewards = &rewards
	return nil
}

// ToBattleResponse converts a BattleResult to the API response format
func (br *BattleResult) ToBattleResponse() map[string]interface{} {
	// Make sure rewards are parsed
	if br.Rewards == nil {
		br.ParseRewardsJSON()
	}
	
	return map[string]interface{}{
		"battle_id":  br.ID,
		"result":     br.Result,
		"battle_log": br.BattleLog,
		"rewards":    br.Rewards,
	}
} 