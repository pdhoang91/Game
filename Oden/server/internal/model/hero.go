package model

import "time"

// HeroType represents a template for heroes
type HeroType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Rarity      string `json:"rarity"` // common, rare, epic, legendary
	BaseHP      int    `json:"base_hp"`
	BaseATK     int    `json:"base_atk"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	Skills      []Skill `json:"skills,omitempty"`
}

// Skill represents a hero ability
type Skill struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description,omitempty"`
	DamageMultiplier float64 `json:"damage_multiplier"`
	Cooldown         int     `json:"cooldown"`
	TargetsAll       bool    `json:"targets_all"`
}

// Hero represents a player's hero
type Hero struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	HeroTypeID string    `json:"hero_type_id"`
	Level      int       `json:"level"`
	Experience int       `json:"experience"`
	CreatedAt  time.Time `json:"created_at"`
	
	// Computed fields (not stored in DB)
	HeroType   *HeroType `json:"hero_type,omitempty"`
	HP         int       `json:"hp,omitempty"`
	ATK        int       `json:"atk,omitempty"`
	Skills     []Skill   `json:"skills,omitempty"`
}

// NewHero creates a new hero instance
func NewHero(id, userID, heroTypeID string) *Hero {
	return &Hero{
		ID:         id,
		UserID:     userID,
		HeroTypeID: heroTypeID,
		Level:      1,
		Experience: 0,
		CreatedAt:  time.Now(),
	}
}

// CalculateStats calculates the hero's stats based on level and base stats
func (h *Hero) CalculateStats() {
	if h.HeroType == nil {
		return
	}
	
	// Simple level-based scaling for MVP
	levelFactor := 1.0 + float64(h.Level-1)*0.1 // 10% increase per level
	
	h.HP = int(float64(h.HeroType.BaseHP) * levelFactor)
	h.ATK = int(float64(h.HeroType.BaseATK) * levelFactor)
}

// AddExperience adds experience to the hero and levels up if necessary
func (h *Hero) AddExperience(amount int) bool {
	oldLevel := h.Level
	h.Experience += amount
	
	// Simple leveling formula: 100 XP per level
	h.Level = 1 + (h.Experience / 100)
	
	// If level changed, recalculate stats
	if h.Level != oldLevel {
		h.CalculateStats()
		return true // Level up occurred
	}
	
	return false // No level up
}

// HeroWithDetails represents a hero with all its details
type HeroWithDetails struct {
	ID         string    `json:"id"`
	HeroTypeID string    `json:"hero_type_id"`
	Name       string    `json:"name"`       // From HeroType
	Level      int       `json:"level"`
	Experience int       `json:"experience"`
	HP         int       `json:"hp"`         // Calculated
	ATK        int       `json:"atk"`        // Calculated
	Skills     []Skill   `json:"skills"`     // From HeroType
}

// ToHeroWithDetails converts a Hero to HeroWithDetails
func (h *Hero) ToHeroWithDetails() *HeroWithDetails {
	if h.HeroType == nil {
		return &HeroWithDetails{
			ID:         h.ID,
			HeroTypeID: h.HeroTypeID,
			Level:      h.Level,
			Experience: h.Experience,
		}
	}
	
	// Calculate stats if not already calculated
	if h.HP == 0 || h.ATK == 0 {
		h.CalculateStats()
	}
	
	return &HeroWithDetails{
		ID:         h.ID,
		HeroTypeID: h.HeroTypeID,
		Name:       h.HeroType.Name,
		Level:      h.Level,
		Experience: h.Experience,
		HP:         h.HP,
		ATK:        h.ATK,
		Skills:     h.Skills,
	}
} 