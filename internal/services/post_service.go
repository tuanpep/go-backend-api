package services

import (
	"time"

	"go-backend-api/internal/models"
	"go-backend-api/internal/pkg/errors"
	"go-backend-api/internal/pkg/validation"

	"github.com/google/uuid"
)

// postService implements PostService interface
type postService struct {
	postRepo  models.PostRepository
	userRepo  models.UserRepository
	validator *validation.Validator
}

// NewPostService creates a new post service
func NewPostService(postRepo models.PostRepository, userRepo models.UserRepository) models.PostService {
	return &postService{
		postRepo:  postRepo,
		userRepo:  userRepo,
		validator: validation.NewValidator(),
	}
}

// CreatePost creates a new post
func (s *postService) CreatePost(authorID uuid.UUID, req *models.CreatePostRequest) (*models.Post, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, errors.WrapErrorWithCode(err, 400, "Validation failed")
	}

	// Verify author exists
	author, err := s.userRepo.GetByID(authorID)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get author")
	}
	if author == nil {
		return nil, errors.ErrUserNotFound
	}

	// Create post
	post := &models.Post{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, errors.WrapError(err, "Failed to create post")
	}

	return post, nil
}

// GetPostByID gets a post by ID
func (s *postService) GetPostByID(id uuid.UUID) (*models.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get post")
	}
	if post == nil {
		return nil, errors.ErrPostNotFound
	}

	// Get author information
	author, err := s.userRepo.GetByID(post.AuthorID)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get post author")
	}
	if author != nil {
		author.Password = "" // Clear password
		post.Author = author
	}

	return post, nil
}

// GetPosts gets all posts with pagination
func (s *postService) GetPosts(page, perPage int) ([]*models.Post, int, error) {
	offset := (page - 1) * perPage

	posts, err := s.postRepo.GetAllWithAuthor(perPage, offset)
	if err != nil {
		return nil, 0, errors.WrapError(err, "Failed to get posts")
	}

	total, err := s.postRepo.Count()
	if err != nil {
		return nil, 0, errors.WrapError(err, "Failed to count posts")
	}

	// Clear passwords from author information
	for _, post := range posts {
		if post.Author != nil {
			post.Author.Password = ""
		}
	}

	return posts, total, nil
}

// GetPostsByAuthor gets posts by author with pagination
func (s *postService) GetPostsByAuthor(authorID uuid.UUID, page, perPage int) ([]*models.Post, int, error) {
	offset := (page - 1) * perPage

	posts, err := s.postRepo.GetByAuthorID(authorID, perPage, offset)
	if err != nil {
		return nil, 0, errors.WrapError(err, "Failed to get posts by author")
	}

	total, err := s.postRepo.CountByAuthorID(authorID)
	if err != nil {
		return nil, 0, errors.WrapError(err, "Failed to count posts by author")
	}

	// Get author information for each post
	for _, post := range posts {
		author, err := s.userRepo.GetByID(post.AuthorID)
		if err != nil {
			return nil, 0, errors.WrapError(err, "Failed to get post author")
		}
		if author != nil {
			author.Password = "" // Clear password
			post.Author = author
		}
	}

	return posts, total, nil
}

// UpdatePost updates a post
func (s *postService) UpdatePost(id, authorID uuid.UUID, req *models.UpdatePostRequest) (*models.Post, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, errors.WrapErrorWithCode(err, 400, "Validation failed")
	}

	// Get existing post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get post")
	}
	if post == nil {
		return nil, errors.ErrPostNotFound
	}

	// Check if user is the author
	if post.AuthorID != authorID {
		return nil, errors.ErrForbidden
	}

	// Update fields if provided
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	post.UpdatedAt = time.Now()

	// Update post
	if err := s.postRepo.Update(post); err != nil {
		return nil, errors.WrapError(err, "Failed to update post")
	}

	// Get author information
	author, err := s.userRepo.GetByID(post.AuthorID)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get post author")
	}
	if author != nil {
		author.Password = "" // Clear password
		post.Author = author
	}

	return post, nil
}

// DeletePost deletes a post
func (s *postService) DeletePost(id, authorID uuid.UUID) error {
	// Get existing post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return errors.WrapError(err, "Failed to get post")
	}
	if post == nil {
		return errors.ErrPostNotFound
	}

	// Check if user is the author
	if post.AuthorID != authorID {
		return errors.ErrForbidden
	}

	// Delete post
	if err := s.postRepo.Delete(id); err != nil {
		return errors.WrapError(err, "Failed to delete post")
	}

	return nil
}

// GetPublishedPosts gets published posts with pagination
func (s *postService) GetPublishedPosts(page, perPage int) ([]*models.Post, int, error) {
	offset := (page - 1) * perPage

	posts, err := s.postRepo.GetPublished(perPage, offset)
	if err != nil {
		return nil, 0, errors.WrapError(err, "Failed to get published posts")
	}

	total, err := s.postRepo.CountPublished()
	if err != nil {
		return nil, 0, errors.WrapError(err, "Failed to count published posts")
	}

	return posts, total, nil
}

// PublishPost publishes a post
func (s *postService) PublishPost(id, authorID uuid.UUID) error {
	// Get existing post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return errors.WrapError(err, "Failed to get post")
	}
	if post == nil {
		return errors.NewErrorWithCode(404, "Post not found")
	}

	// Check if user owns the post
	if post.AuthorID != authorID {
		return errors.NewErrorWithCode(403, "Not authorized to publish this post")
	}

	// Update post
	post.IsPublished = true
	post.UpdatedAt = time.Now()

	err = s.postRepo.Update(post)
	if err != nil {
		return errors.WrapError(err, "Failed to publish post")
	}

	return nil
}

// UnpublishPost unpublishes a post
func (s *postService) UnpublishPost(id, authorID uuid.UUID) error {
	// Get existing post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return errors.WrapError(err, "Failed to get post")
	}
	if post == nil {
		return errors.NewErrorWithCode(404, "Post not found")
	}

	// Check if user owns the post
	if post.AuthorID != authorID {
		return errors.NewErrorWithCode(403, "Not authorized to unpublish this post")
	}

	// Update post
	post.IsPublished = false
	post.UpdatedAt = time.Now()

	err = s.postRepo.Update(post)
	if err != nil {
		return errors.WrapError(err, "Failed to unpublish post")
	}

	return nil
}

// ValidatePost validates a post entity
func (s *postService) ValidatePost(post *models.Post) error {
	return s.validator.Validate(post)
}
