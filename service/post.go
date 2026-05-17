package service

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/repository"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type IPostService interface {
	Create(c *gin.Context, post *model.Post) error
	GetTimeline(c *gin.Context) ([]model.Post, error)
	GetTimelineNoCache() ([]model.Post, error)
	Delete(c *gin.Context, userID, postID uint) error
}

type PostService struct {
	repo      repository.IPostRepository
	uploadDir string
	maxSize   int64
}

func NewPostService(repo repository.IPostRepository, uploadDir string, maxSize int64) IPostService {
	return &PostService{repo: repo, uploadDir: uploadDir, maxSize: maxSize}
}

func (s *PostService) Create(c *gin.Context, post *model.Post) error {
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["images"]
		var filenames []string

		for _, file := range files {
			if file.Size > s.maxSize {
				return fmt.Errorf("图片大小超过限制（最大 %d 字节）", s.maxSize)
			}

			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
				return fmt.Errorf("不支持的图片格式: %s", ext)
			}

			filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			savePath := filepath.Join(s.uploadDir, filename)

			if err := c.SaveUploadedFile(file, savePath); err != nil {
				return fmt.Errorf("保存图片失败: %w", err)
			}

			filenames = append(filenames, filename)
		}

		if len(filenames) > 0 {
			post.Images = strings.Join(filenames, ",")
		}
	}

	// 清除时间线缓存
	ctx := c.Request.Context()
	if err := global.Redis.Del(ctx, "posts:timeline").Err(); err != nil {
		global.Logger.Warnw("清除时间线缓存失败", "error", err)
	}

	return s.repo.Create(post)
}

func (s *PostService) GetTimeline(c *gin.Context) ([]model.Post, error) {
	ctx := c.Request.Context()

	if str := global.Redis.Get(ctx, "posts:timeline").Val(); str != "" {
		var posts []model.Post
		if err := json.Unmarshal([]byte(str), &posts); err == nil {
			return posts, nil
		}
	}

	posts, err := s.repo.GetTimeline()
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(posts)
	if err != nil {
		return nil, err
	}
	global.Redis.Set(ctx, "posts:timeline", jsonData, time.Minute*10)

	return posts, nil
}

func (s *PostService) GetTimelineNoCache() ([]model.Post, error) {
	return s.repo.GetTimeline()
}

func (s *PostService) Delete(c *gin.Context, userID, postID uint) error {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return fmt.Errorf("动态不存在")
	}

	if post.UserID != userID {
		return fmt.Errorf("无权限删除该动态")
	}

	// 删除关联图片
	if post.Images != "" {
		for _, filename := range strings.Split(post.Images, ",") {
			filename = strings.TrimSpace(filename)
			if filename != "" {
				os.Remove(filepath.Join(s.uploadDir, filename))
			}
		}
	}

	// 清除时间线缓存
	ctx := c.Request.Context()
	if err := global.Redis.Del(ctx, "posts:timeline").Err(); err != nil {
		global.Logger.Warnw("清除时间线缓存失败", "error", err)
	}

	return s.repo.Delete(postID)
}
