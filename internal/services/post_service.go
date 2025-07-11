package services

import (
	"gorm.io/gorm"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
)

type PostService struct {
	BaseService
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		BaseService: BaseService{DB: db},
	}
}

func (s *PostService) CreatePost(req *models.PostRequest) (*models.Post, error) {
	post := &models.Post{
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
	}

	err := s.Create(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostByID(id uint64) (*models.Post, error) {
	var post models.Post
	err := s.GetByID(&post, id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetAllPosts(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := s.GetAll(&posts, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) UpdatePost(id uint64, req *models.PostRequest) (*models.Post, error) {
	updates := map[string]interface{}{
		"title":      req.Title,
		"content":    req.Content,
		"created_by": req.CreatedBy,
	}

	err := s.Update(&models.Post{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetPostByID(id)
}

func (s *PostService) DeletePost(id uint64) error {
	return s.Delete(&models.Post{}, id)
}

func (s *PostService) GetPostsByAuthor(authorID uint64, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	query := s.DB.Where("created_by = ?", authorID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
