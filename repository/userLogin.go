package repository

import (
	"IndulgenceMealPlan/model"

	"gorm.io/gorm"
)

type LoginRepository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) *LoginRepository {
	return &LoginRepository{db: db}
}

func (r *LoginRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *LoginRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *LoginRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *LoginRepository) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
