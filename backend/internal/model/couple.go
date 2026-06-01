package model

import "time"

type CoupleBinding struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserAID    uint      `gorm:"not null;index" json:"user_a_id"`
	UserBID    *uint     `gorm:"index" json:"user_b_id"`
	CoupleName string    `gorm:"type:varchar(50)" json:"couple_name"`
	InviteCode string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"invite_code"`
	Status     int8      `gorm:"default:1" json:"status"` // 0-已解绑 1-已绑定
	BoundAt    time.Time `json:"bound_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// relations
	UserA *User `gorm:"foreignKey:UserAID" json:"user_a,omitempty"`
	UserB *User `gorm:"foreignKey:UserBID" json:"user_b,omitempty"`
}

func (CoupleBinding) TableName() string { return "couple_bindings" }

type CoupleOrder struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CoupleID  uint      `gorm:"not null;index" json:"couple_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	RecipeID  *uint     `json:"recipe_id"` // nullable
	DishName  string    `gorm:"type:varchar(100);not null" json:"dish_name"`
	MealType  string    `gorm:"type:varchar(20)" json:"meal_type"`
	MealDate  string    `gorm:"type:varchar(10)" json:"meal_date"` // YYYY-MM-DD
	Note      string    `gorm:"type:varchar(200)" json:"note"`
	Status    int8      `gorm:"default:0" json:"status"` // 0-待确认 1-已确认 2-已取消
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// relations
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Recipe *Recipe `gorm:"foreignKey:RecipeID" json:"recipe,omitempty"`
}

func (CoupleOrder) TableName() string { return "couple_orders" }
