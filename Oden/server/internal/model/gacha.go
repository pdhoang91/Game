package model

import (
	"encoding/json"
	"time"
)

// BannerType represents the type of summon banner
type BannerType string

const (
	BannerTypeStandard  BannerType = "standard"  // Always available
	BannerTypeEvent     BannerType = "event"     // Limited time
	BannerTypeSpecial   BannerType = "special"   // Special banners (e.g., beginner, guaranteed)
)

// SummonCostType represents the type of currency used for summons
type SummonCostType string

const (
	SummonCostGem        SummonCostType = "gem"
	SummonCostSummonTicket SummonCostType = "summon_ticket"
	SummonCostSpecialTicket SummonCostType = "special_ticket"
)

// Banner represents a summon banner
type Banner struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        BannerType `json:"type"`
	ImageURL    string     `json:"image_url"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"` // Null for standard banners
	
	// Rates
	StandardHeroRate   float64 `json:"standard_hero_rate"`     // Base rate for 5* heroes
	FeaturedHeroRate   float64 `json:"featured_hero_rate"`     // Rate for featured heroes
	GuaranteeThreshold int     `json:"guarantee_threshold"`    // Guaranteed 5* after this many pulls
	
	// Cost
	SingleSummonCost    int             `json:"single_summon_cost"`
	TenSummonCost       int             `json:"ten_summon_cost"`
	CostType            SummonCostType  `json:"cost_type"`
	
	// Featured heroes/items
	FeaturedHeroes []string `json:"featured_heroes"` // HeroType IDs
	FeaturedItems  []string `json:"featured_items"`  // ItemTemplate IDs
	
	// Pool of available heroes/items (if empty, use standard pool)
	HeroPool       []string `json:"hero_pool,omitempty"`
	ItemPool       []string `json:"item_pool,omitempty"`
	
	// Daily free summon
	HasDailyFreeSummon bool `json:"has_daily_free_summon"`
}

// IsActive checks if the banner is currently active
func (b *Banner) IsActive() bool {
	now := time.Now()
	if now.Before(b.StartTime) {
		return false
	}
	
	if b.EndTime != nil && now.After(*b.EndTime) {
		return false
	}
	
	return true
}

// SummonResult represents the result of a summon
type SummonResult struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	BannerID  string    `json:"banner_id"`
	Timestamp time.Time `json:"timestamp"`
	
	// Results
	ResultType   string `json:"result_type"` // hero or item
	ResultID     string `json:"result_id"`   // HeroType ID or ItemTemplate ID
	Rarity       string `json:"rarity"`      // common, uncommon, rare, epic, legendary
	IsFeatured   bool   `json:"is_featured"`
	IsPityBreak  bool   `json:"is_pity_break"`
	
	// For tracking pity system
	PullNumber   int    `json:"pull_number"`
}

// NewSummonResult creates a new summon result
func NewSummonResult(id, userID, bannerID, resultType, resultID, rarity string, isFeatured, isPityBreak bool, pullNumber int) *SummonResult {
	return &SummonResult{
		ID:         id,
		UserID:     userID,
		BannerID:   bannerID,
		Timestamp:  time.Now(),
		ResultType: resultType,
		ResultID:   resultID,
		Rarity:     rarity,
		IsFeatured: isFeatured,
		IsPityBreak: isPityBreak,
		PullNumber: pullNumber,
	}
}

// SummonSession tracks a user's summon session
type SummonSession struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	BannerID   string    `json:"banner_id"`
	PullCount  int       `json:"pull_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Pity system tracking
	LastLegendaryAt int       `json:"last_legendary_at"` // Pull number of last legendary
	HasGuarantee    bool      `json:"has_guarantee"`     // Next 5* is guaranteed featured
	
	// Daily free summon tracking
	LastFreeSummon  *time.Time `json:"last_free_summon,omitempty"`
}

// NewSummonSession creates a new summon session
func NewSummonSession(id, userID, bannerID string) *SummonSession {
	now := time.Now()
	return &SummonSession{
		ID:              id,
		UserID:          userID,
		BannerID:        bannerID,
		PullCount:       0,
		CreatedAt:       now,
		UpdatedAt:       now,
		LastLegendaryAt: 0,
		HasGuarantee:    false,
	}
}

// CanClaimFreeSummon checks if the user can claim a free summon
func (s *SummonSession) CanClaimFreeSummon() bool {
	if s.LastFreeSummon == nil {
		return true
	}
	
	lastFreeDay := s.LastFreeSummon.Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)
	
	return today.After(lastFreeDay)
}

// UpdateFreeSummon updates the last free summon time
func (s *SummonSession) UpdateFreeSummon() {
	now := time.Now()
	s.LastFreeSummon = &now
	s.UpdatedAt = now
}

// IncrementPullCount increments the pull count
func (s *SummonSession) IncrementPullCount() {
	s.PullCount++
	s.UpdatedAt = time.Now()
}

// SummonMultiResult represents the result of a multi-summon
type SummonMultiResult struct {
	BannerID    string          `json:"banner_id"`
	BannerName  string          `json:"banner_name"`
	Results     []*SummonResult `json:"results"`
	NewHeroes   []*HeroWithDetails  `json:"new_heroes,omitempty"`
	NewItems    []*ItemWithTemplate `json:"new_items,omitempty"`
}

// SummonRateInfo represents summon rate information for the client
type SummonRateInfo struct {
	BannerID            string  `json:"banner_id"`
	LegendaryRate       float64 `json:"legendary_rate"`
	FeaturedHeroRate    float64 `json:"featured_hero_rate"`
	GuaranteeThreshold  int     `json:"guarantee_threshold"`
	CurrentPity         int     `json:"current_pity"`
	HasGuaranteeActive  bool    `json:"has_guarantee_active"`
	FeaturedHeroes      []HeroBasicInfo `json:"featured_heroes"`
}

// HeroBasicInfo represents basic hero information for display in summon rates
type HeroBasicInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Rarity      string     `json:"rarity"`
	ImageURL    string     `json:"image_url,omitempty"`
}

// ToSummonRateInfo converts a Banner and SummonSession to SummonRateInfo
func (b *Banner) ToSummonRateInfo(session *SummonSession, featuredHeroes []*HeroType) *SummonRateInfo {
	heroInfos := make([]HeroBasicInfo, 0, len(featuredHeroes))
	for _, hero := range featuredHeroes {
		heroInfos = append(heroInfos, HeroBasicInfo{
			ID:       hero.ID,
			Name:     hero.Name,
			Rarity:   string(hero.Rarity),
			ImageURL: hero.ImageURL,
		})
	}
	
	currentPity := 0
	hasGuarantee := false
	
	if session != nil {
		if session.LastLegendaryAt > 0 {
			currentPity = session.PullCount - session.LastLegendaryAt
		} else {
			currentPity = session.PullCount
		}
		hasGuarantee = session.HasGuarantee
	}
	
	return &SummonRateInfo{
		BannerID:           b.ID,
		LegendaryRate:      b.StandardHeroRate + b.FeaturedHeroRate,
		FeaturedHeroRate:   b.FeaturedHeroRate,
		GuaranteeThreshold: b.GuaranteeThreshold,
		CurrentPity:        currentPity,
		HasGuaranteeActive: hasGuarantee,
		FeaturedHeroes:     heroInfos,
	}
} 