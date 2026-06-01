package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type CoupleRepo struct {
	db *gorm.DB
}

func NewCoupleRepo(db *gorm.DB) *CoupleRepo {
	return &CoupleRepo{db: db}
}

// ── CoupleBinding ──

func (r *CoupleRepo) CreateBinding(b *model.CoupleBinding) error {
	// 使用 Exec 绕过 GORM 零值跳过问题，确保 Status=0 被写入
	err := r.db.Exec(
		"INSERT INTO couple_bindings (user_a_id, couple_name, invite_code, status, bound_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
		b.UserAID, b.CoupleName, b.InviteCode, b.Status, b.BoundAt,
	).Error
	if err != nil {
		return err
	}
	// 查回刚创建的记录获取 ID
	return r.db.Where("invite_code = ?", b.InviteCode).Order("id DESC").First(b).Error
}

func (r *CoupleRepo) FindBindingByInviteCode(code string) (*model.CoupleBinding, error) {
	var b model.CoupleBinding
	err := r.db.Where("invite_code = ? AND status = 0", code).First(&b).Error
	return &b, err
}

func (r *CoupleRepo) FindActiveBindingByUser(userID uint) (*model.CoupleBinding, error) {
	var b model.CoupleBinding
	err := r.db.Where("(user_a_id = ? OR user_b_id = ?) AND status = 1", userID, userID).
		Preload("UserA").Preload("UserB").First(&b).Error
	return &b, err
}

func (r *CoupleRepo) UpdateBinding(b *model.CoupleBinding) error {
	return r.db.Save(b).Error
}

// CleanupStaleInvites 删除用户旧的未完成邀请（status=0）
func (r *CoupleRepo) CleanupStaleInvites(userID uint) error {
	return r.db.Where("user_a_id = ? AND status = 0", userID).Delete(&model.CoupleBinding{}).Error
}

// ── CoupleOrder ──

func (r *CoupleRepo) CreateOrder(o *model.CoupleOrder) error {
	return r.db.Create(o).Error
}

func (r *CoupleRepo) FindOrdersByCoupleID(coupleID uint, mealDate string) ([]model.CoupleOrder, error) {
	var orders []model.CoupleOrder
	query := r.db.Where("couple_id = ?", coupleID).Preload("User")
	if mealDate != "" {
		query = query.Where("meal_date = ?", mealDate)
	}
	err := query.Order("meal_date DESC, created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *CoupleRepo) FindOrdersByCoupleIDAndType(coupleID uint, mealDate, mealType string) ([]model.CoupleOrder, error) {
	var orders []model.CoupleOrder
	query := r.db.Where("couple_id = ? AND status = 1", coupleID)
	if mealDate != "" {
		query = query.Where("meal_date = ?", mealDate)
	}
	if mealType != "" {
		query = query.Where("meal_type = ?", mealType)
	}
	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *CoupleRepo) FindOrderByID(id uint) (*model.CoupleOrder, error) {
	var o model.CoupleOrder
	err := r.db.First(&o, id).Error
	return &o, err
}

func (r *CoupleRepo) UpdateOrder(o *model.CoupleOrder) error {
	return r.db.Save(o).Error
}

func (r *CoupleRepo) DeleteOrder(id uint) error {
	return r.db.Delete(&model.CoupleOrder{}, id).Error
}
