package model

import (
	"encoding/json"
	"time"
)

// MissionType represents the type of mission
type MissionType string

const (
	MissionTypeDaily   MissionType = "daily"
	MissionTypeWeekly  MissionType = "weekly"
	MissionTypeStory   MissionType = "story"
	MissionTypeAchievement MissionType = "achievement"
)

// MissionStatus represents the status of a mission
type MissionStatus string

const (
	MissionStatusInProgress MissionStatus = "in_progress"
	MissionStatusCompleted  MissionStatus = "completed"
	MissionStatusClaimed    MissionStatus = "claimed"
)

// MissionRequirementType represents the type of requirement for a mission
type MissionRequirementType string

const (
	// Battle related
	RequirementCompleteBattles MissionRequirementType = "complete_battles"
	RequirementWinBattles      MissionRequirementType = "win_battles"
	RequirementKillEnemies     MissionRequirementType = "kill_enemies"
	
	// Hero related
	RequirementLevelUpHero     MissionRequirementType = "level_up_hero"
	RequirementOwnHeroes       MissionRequirementType = "own_heroes"
	RequirementMaxLevelHero    MissionRequirementType = "max_level_hero"
	
	// Item related
	RequirementCollectItems    MissionRequirementType = "collect_items"
	RequirementEquipItems      MissionRequirementType = "equip_items"
	RequirementUpgradeItems    MissionRequirementType = "upgrade_items"
	
	// Resource related
	RequirementSpendGold       MissionRequirementType = "spend_gold"
	RequirementSpendGems       MissionRequirementType = "spend_gems"
)

// MissionTemplate represents a template for a mission
type MissionTemplate struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Type        MissionType `json:"type"`
	
	// Requirements to complete the mission
	RequirementType MissionRequirementType `json:"requirement_type"`
	TargetValue     int                   `json:"target_value"`
	TargetID        string                `json:"target_id,omitempty"` // Optional ID for specific target (e.g., specific hero/item)
	
	// Rewards
	GoldReward      int      `json:"gold_reward"`
	GemsReward      int      `json:"gems_reward"`
	ExperienceReward int      `json:"experience_reward"`
	ItemRewards     []string `json:"item_rewards,omitempty"` // ItemTemplate IDs
}

// Mission represents a mission assigned to a user
type Mission struct {
	ID              string        `json:"id"`
	UserID          string        `json:"user_id"`
	MissionTemplateID string      `json:"mission_template_id"`
	Status          MissionStatus `json:"status"`
	CurrentValue    int           `json:"current_value"`
	AssignedAt      time.Time     `json:"assigned_at"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
	ClaimedAt       *time.Time    `json:"claimed_at,omitempty"`
	ExpiresAt       *time.Time    `json:"expires_at,omitempty"` // For daily/weekly missions
	
	// Computed fields
	Template        *MissionTemplate `json:"template,omitempty"`
}

// NewMission creates a new mission instance
func NewMission(id, userID, missionTemplateID string, template *MissionTemplate) *Mission {
	mission := &Mission{
		ID:                id,
		UserID:            userID,
		MissionTemplateID: missionTemplateID,
		Status:            MissionStatusInProgress,
		CurrentValue:      0,
		AssignedAt:        time.Now(),
		Template:          template,
	}
	
	// Set expiration for time-limited missions
	if template != nil {
		switch template.Type {
		case MissionTypeDaily:
			expiry := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
			mission.ExpiresAt = &expiry
		case MissionTypeWeekly:
			expiry := time.Now().Add(7 * 24 * time.Hour).Truncate(24 * time.Hour)
			mission.ExpiresAt = &expiry
		}
	}
	
	return mission
}

// UpdateProgress updates the mission progress
func (m *Mission) UpdateProgress(value int) bool {
	if m.Status != MissionStatusInProgress {
		return false
	}
	
	m.CurrentValue += value
	
	if m.Template != nil && m.CurrentValue >= m.Template.TargetValue {
		m.CurrentValue = m.Template.TargetValue
		m.Status = MissionStatusCompleted
		now := time.Now()
		m.CompletedAt = &now
		return true // Mission completed
	}
	
	return false // Mission still in progress
}

// ClaimRewards marks the mission as claimed
func (m *Mission) ClaimRewards() bool {
	if m.Status != MissionStatusCompleted {
		return false
	}
	
	m.Status = MissionStatusClaimed
	now := time.Now()
	m.ClaimedAt = &now
	return true
}

// IsExpired checks if the mission is expired
func (m *Mission) IsExpired() bool {
	if m.ExpiresAt == nil {
		return false
	}
	
	return time.Now().After(*m.ExpiresAt)
}

// MissionProgress represents mission progress for the client
type MissionProgress struct {
	MissionID    string        `json:"mission_id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Type         MissionType   `json:"type"`
	Status       MissionStatus `json:"status"`
	CurrentValue int           `json:"current_value"`
	TargetValue  int           `json:"target_value"`
	Rewards      MissionRewards `json:"rewards"`
	ExpiresAt    *time.Time    `json:"expires_at,omitempty"`
}

// MissionRewards represents the rewards for a mission
type MissionRewards struct {
	Gold       int           `json:"gold"`
	Gems       int           `json:"gems"`
	Experience int           `json:"experience"`
	Items      []ItemReward  `json:"items,omitempty"`
}

// ItemReward represents an item reward for a mission
type ItemReward struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// ToMissionProgress converts a Mission to MissionProgress
func (m *Mission) ToMissionProgress() *MissionProgress {
	if m.Template == nil {
		return &MissionProgress{
			MissionID:    m.ID,
			Status:       m.Status,
			CurrentValue: m.CurrentValue,
			ExpiresAt:    m.ExpiresAt,
		}
	}
	
	return &MissionProgress{
		MissionID:    m.ID,
		Title:        m.Template.Title,
		Description:  m.Template.Description,
		Type:         m.Template.Type,
		Status:       m.Status,
		CurrentValue: m.CurrentValue,
		TargetValue:  m.Template.TargetValue,
		Rewards: MissionRewards{
			Gold:       m.Template.GoldReward,
			Gems:       m.Template.GemsReward,
			Experience: m.Template.ExperienceReward,
			// Items will be populated by the service layer
		},
		ExpiresAt:    m.ExpiresAt,
	}
}