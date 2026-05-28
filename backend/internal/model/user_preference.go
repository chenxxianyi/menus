package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type UserPreference struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	TastePreference     JSON      `gorm:"type:json" json:"taste_preference"`
	HealthGoal          string    `gorm:"type:varchar(50)" json:"health_goal"`
	AvoidIngredients    JSON      `gorm:"type:json" json:"avoid_ingredients"`
	FavoriteIngredients JSON      `gorm:"type:json" json:"favorite_ingredients"`
	CookTimePreference  string    `gorm:"type:varchar(50)" json:"cook_time_preference"`
	PeopleCount         int       `gorm:"default:2" json:"people_count"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (UserPreference) TableName() string { return "user_preferences" }

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = append((*j)[0:0], []byte(v)...)
	default:
		return fmt.Errorf("failed to scan JSON: %v", value)
	}
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return fmt.Errorf("JSON UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
