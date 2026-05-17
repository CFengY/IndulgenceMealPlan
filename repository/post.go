package repository

import (
	"IndulgenceMealPlan/model"

	"gorm.io/gorm"
)

type IPostRepository interface {
	Create(post *model.Post) error
	GetByID(id uint) (*model.Post, error)
	GetTimeline() ([]model.Post, error)
	Delete(id uint) error
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetTimeline() ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Preload("User").Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}
