package models

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post entity
type Post struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	AuthorID    uuid.UUID `json:"author_id" db:"author_id"`
	Author      *User     `json:"author,omitempty" db:"-"`
	IsPublished bool      `json:"is_published" db:"is_published"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PostRepository defines the interface for post data operations
type PostRepository interface {
	Create(post *Post) error
	GetByID(id uuid.UUID) (*Post, error)
	GetByAuthorID(authorID uuid.UUID, limit, offset int) ([]*Post, error)
	GetAll(limit, offset int) ([]*Post, error)
	GetAllWithAuthor(limit, offset int) ([]*Post, error)
	GetPublished(limit, offset int) ([]*Post, error)
	Update(post *Post) error
	Delete(id uuid.UUID) error
	Count() (int, error)
	CountByAuthorID(authorID uuid.UUID) (int, error)
	CountPublished() (int, error)
}

// PostService defines the interface for post business logic
type PostService interface {
	CreatePost(authorID uuid.UUID, req *CreatePostRequest) (*Post, error)
	GetPostByID(id uuid.UUID) (*Post, error)
	GetPosts(page, perPage int) ([]*Post, int, error)
	GetPostsByAuthor(authorID uuid.UUID, page, perPage int) ([]*Post, int, error)
	GetPublishedPosts(page, perPage int) ([]*Post, int, error)
	UpdatePost(id, authorID uuid.UUID, req *UpdatePostRequest) (*Post, error)
	DeletePost(id, authorID uuid.UUID) error
	PublishPost(id, authorID uuid.UUID) error
	UnpublishPost(id, authorID uuid.UUID) error
	ValidatePost(post *Post) error
}

// CreatePostRequest represents the request to create a post
type CreatePostRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Content     string `json:"content" validate:"required,min=1"`
	IsPublished bool   `json:"is_published,omitempty"`
}

// UpdatePostRequest represents the request to update a post
type UpdatePostRequest struct {
	Title       string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Content     string `json:"content,omitempty" validate:"omitempty,min=1"`
	IsPublished *bool  `json:"is_published,omitempty"`
}

// PostWithAuthor represents a post with author information
type PostWithAuthor struct {
	Post
	Author User `json:"author"`
}
