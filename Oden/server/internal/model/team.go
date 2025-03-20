package model

import "time"

// Team represents a player's team of heroes
type Team struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Position1  string    `json:"position_1,omitempty"`
	Position2  string    `json:"position_2,omitempty"`
	Position3  string    `json:"position_3,omitempty"`
	Position4  string    `json:"position_4,omitempty"`
	Position5  string    `json:"position_5,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Computed fields (not stored in DB)
	Heroes     map[int]*Hero `json:"heroes,omitempty"`
}

// NewTeam creates a new team instance
func NewTeam(id, userID string) *Team {
	now := time.Now()
	return &Team{
		ID:        id,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		Heroes:    make(map[int]*Hero),
	}
}

// GetPositionHeroID returns the hero ID at the specified position
func (t *Team) GetPositionHeroID(position int) string {
	switch position {
	case 1:
		return t.Position1
	case 2:
		return t.Position2
	case 3:
		return t.Position3
	case 4:
		return t.Position4
	case 5:
		return t.Position5
	default:
		return ""
	}
}

// SetPositionHeroID sets the hero ID at the specified position
func (t *Team) SetPositionHeroID(position int, heroID string) {
	switch position {
	case 1:
		t.Position1 = heroID
	case 2:
		t.Position2 = heroID
	case 3:
		t.Position3 = heroID
	case 4:
		t.Position4 = heroID
	case 5:
		t.Position5 = heroID
	}
	t.UpdatedAt = time.Now()
}

// GetAllPositions returns a map of position to hero ID
func (t *Team) GetAllPositions() map[int]string {
	positions := make(map[int]string)
	
	if t.Position1 != "" {
		positions[1] = t.Position1
	}
	if t.Position2 != "" {
		positions[2] = t.Position2
	}
	if t.Position3 != "" {
		positions[3] = t.Position3
	}
	if t.Position4 != "" {
		positions[4] = t.Position4
	}
	if t.Position5 != "" {
		positions[5] = t.Position5
	}
	
	return positions
}

// SetAllPositions sets all positions from a map
func (t *Team) SetAllPositions(positions map[string]string) {
	// Clear all positions first
	t.Position1 = ""
	t.Position2 = ""
	t.Position3 = ""
	t.Position4 = ""
	t.Position5 = ""
	
	// Set positions from map
	for posStr, heroID := range positions {
		var pos int
		switch posStr {
		case "1":
			pos = 1
		case "2":
			pos = 2
		case "3":
			pos = 3
		case "4":
			pos = 4
		case "5":
			pos = 5
		default:
			continue // Skip invalid positions
		}
		
		if heroID != "" && heroID != "null" {
			t.SetPositionHeroID(pos, heroID)
		}
	}
}

// CountHeroes returns the number of heroes in the team
func (t *Team) CountHeroes() int {
	count := 0
	if t.Position1 != "" {
		count++
	}
	if t.Position2 != "" {
		count++
	}
	if t.Position3 != "" {
		count++
	}
	if t.Position4 != "" {
		count++
	}
	if t.Position5 != "" {
		count++
	}
	return count
}

// GetHeroIDs returns all hero IDs in the team
func (t *Team) GetHeroIDs() []string {
	var ids []string
	if t.Position1 != "" {
		ids = append(ids, t.Position1)
	}
	if t.Position2 != "" {
		ids = append(ids, t.Position2)
	}
	if t.Position3 != "" {
		ids = append(ids, t.Position3)
	}
	if t.Position4 != "" {
		ids = append(ids, t.Position4)
	}
	if t.Position5 != "" {
		ids = append(ids, t.Position5)
	}
	return ids
}

// TeamResponse represents the team data sent to the client
type TeamResponse struct {
	TeamID    string                   `json:"team_id"`
	Positions map[string]*HeroPosition `json:"positions"`
}

// HeroPosition represents a hero in a team position
type HeroPosition struct {
	HeroID     string `json:"hero_id"`
	HeroTypeID string `json:"hero_type_id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
}

// ToTeamResponse converts a Team to TeamResponse
func (t *Team) ToTeamResponse() *TeamResponse {
	res := &TeamResponse{
		TeamID:    t.ID,
		Positions: make(map[string]*HeroPosition),
	}
	
	// Add heroes to positions
	if t.Heroes != nil {
		for pos, hero := range t.Heroes {
			if hero != nil && hero.HeroType != nil {
				res.Positions[string(rune('0'+pos))] = &HeroPosition{
					HeroID:     hero.ID,
					HeroTypeID: hero.HeroTypeID,
					Name:       hero.HeroType.Name,
					Level:      hero.Level,
				}
			}
		}
	} else {
		// If heroes not loaded, just set IDs
		positions := t.GetAllPositions()
		for pos, heroID := range positions {
			res.Positions[string(rune('0'+pos))] = &HeroPosition{
				HeroID: heroID,
			}
		}
	}
	
	return res
} 