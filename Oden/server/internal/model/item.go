package model

import (
	"time"
)

// ItemType represents the category of an item
type ItemType string

const (
	ItemTypeEquipment  ItemType = "equipment"
	ItemTypeConsumable ItemType = "consumable"
	ItemTypeMaterial   ItemType = "material"
)

// ItemRarity represents the rarity of an item
type ItemRarity string

const (
	ItemRarityCommon    ItemRarity = "common"
	ItemRarityUncommon  ItemRarity = "uncommon"
	ItemRarityRare      ItemRarity = "rare"
	ItemRarityEpic      ItemRarity = "epic"
	ItemRarityLegendary ItemRarity = "legendary"
)

// Equipment slots
type EquipmentSlot string

const (
	EquipmentSlotWeapon  EquipmentSlot = "weapon"
	EquipmentSlotArmor   EquipmentSlot = "armor"
	EquipmentSlotAccessory EquipmentSlot = "accessory"
)

// ItemTemplate represents a template for item types in the game
type ItemTemplate struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        ItemType   `json:"type"`
	Rarity      ItemRarity `json:"rarity"`
	ImageURL    string     `json:"image_url,omitempty"`
	
	// Equipment specific
	Slot         EquipmentSlot `json:"slot,omitempty"`
	ATKBonus     int           `json:"atk_bonus,omitempty"`
	HPBonus      int           `json:"hp_bonus,omitempty"`
	
	// Consumable specific
	Effect       string     `json:"effect,omitempty"`
	EffectValue  int        `json:"effect_value,omitempty"`
	
	// Material specific
	UsedForCrafting []string `json:"used_for_crafting,omitempty"`
}

// Item represents a specific instance of an item owned by a player
type Item struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ItemTemplateID string  `json:"item_template_id"`
	Quantity     int       `json:"quantity"`
	AcquiredAt   time.Time `json:"acquired_at"`
	
	// For equipped items
	EquippedToHeroID string `json:"equipped_to_hero_id,omitempty"`
	
	// Computed fields
	Template *ItemTemplate `json:"template,omitempty"`
}

// NewItem creates a new item instance
func NewItem(id, userID, itemTemplateID string, quantity int) *Item {
	return &Item{
		ID:             id,
		UserID:         userID,
		ItemTemplateID: itemTemplateID,
		Quantity:       quantity,
		AcquiredAt:     time.Now(),
	}
}

// IsEquipment checks if the item is equipment
func (i *Item) IsEquipment() bool {
	if i.Template != nil {
		return i.Template.Type == ItemTypeEquipment
	}
	return false
}

// IsConsumable checks if the item is a consumable
func (i *Item) IsConsumable() bool {
	if i.Template != nil {
		return i.Template.Type == ItemTypeConsumable
	}
	return false
}

// IsMaterial checks if the item is a material
func (i *Item) IsMaterial() bool {
	if i.Template != nil {
		return i.Template.Type == ItemTypeMaterial
	}
	return false
}

// EquipToHero equips the item to a hero
func (i *Item) EquipToHero(heroID string) error {
	if !i.IsEquipment() {
		return ErrNotEquipment
	}
	
	i.EquippedToHeroID = heroID
	return nil
}

// UnequipFromHero unequips the item from a hero
func (i *Item) UnequipFromHero() {
	i.EquippedToHeroID = ""
}

// Errors for item operations
var (
	ErrNotEquipment = CustomError{Message: "item is not equipment", Code: "invalid_item_type"}
)

// CustomError represents a custom error with message and code
type CustomError struct {
	Message string
	Code    string
}

// Error implements the error interface
func (e CustomError) Error() string {
	return e.Message
}

// ItemWithTemplate represents an item with its template data
type ItemWithTemplate struct {
	Item
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        ItemType   `json:"type"`
	Rarity      ItemRarity `json:"rarity"`
	ImageURL    string     `json:"image_url,omitempty"`
	Slot        EquipmentSlot `json:"slot,omitempty"`
	ATKBonus    int        `json:"atk_bonus,omitempty"`
	HPBonus     int        `json:"hp_bonus,omitempty"`
	Effect      string     `json:"effect,omitempty"`
	EffectValue int        `json:"effect_value,omitempty"`
}

// ToItemWithTemplate converts an Item to ItemWithTemplate
func (i *Item) ToItemWithTemplate() *ItemWithTemplate {
	if i.Template == nil {
		return &ItemWithTemplate{
			Item: *i,
		}
	}
	
	return &ItemWithTemplate{
		Item:        *i,
		Name:        i.Template.Name,
		Description: i.Template.Description,
		Type:        i.Template.Type,
		Rarity:      i.Template.Rarity,
		ImageURL:    i.Template.ImageURL,
		Slot:        i.Template.Slot,
		ATKBonus:    i.Template.ATKBonus,
		HPBonus:     i.Template.HPBonus,
		Effect:      i.Template.Effect,
		EffectValue: i.Template.EffectValue,
	}
} 