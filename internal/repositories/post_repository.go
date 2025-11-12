package repositories

import (
	"database/sql"

	"go-backend-api/internal/models"
	"go-backend-api/internal/pkg/errors"

	"github.com/google/uuid"
)

// postRepository implements PostRepository interface
type postRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) models.PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post
func (r *postRepository) Create(post *models.Post) error {
	query := `INSERT INTO posts (title, content, author_id, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRow(query, post.Title, post.Content, post.AuthorID, post.CreatedAt, post.UpdatedAt).Scan(&post.ID)
	if err != nil {
		return errors.WrapError(err, "Failed to create post")
	}

	return nil
}

// GetByID gets a post by ID
func (r *postRepository) GetByID(id uuid.UUID) (*models.Post, error) {
	post := &models.Post{}
	query := `SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to get post by ID")
	}

	return post, nil
}

// GetByAuthorID gets posts by author ID
func (r *postRepository) GetByAuthorID(authorID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
			  FROM posts WHERE author_id = $1 
			  ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get posts by author ID")
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to scan post")
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetAll gets all posts
func (r *postRepository) GetAll(limit, offset int) ([]*models.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
			  FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get all posts")
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to scan post")
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetAllWithAuthor gets all posts with author information
func (r *postRepository) GetAllWithAuthor(limit, offset int) ([]*models.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at,
			  u.id, u.username, u.email, u.created_at, u.updated_at
			  FROM posts p
			  JOIN users u ON p.author_id = u.id
			  ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get posts with author")
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		author := &models.User{}

		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to scan post with author")
		}

		post.Author = author
		posts = append(posts, post)
	}

	return posts, nil
}

// Update updates a post
func (r *postRepository) Update(post *models.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, is_published = $3, updated_at = $4 WHERE id = $5`

	_, err := r.db.Exec(query, post.Title, post.Content, post.IsPublished, post.UpdatedAt, post.ID)
	if err != nil {
		return errors.WrapError(err, "Failed to update post")
	}

	return nil
}

// Delete deletes a post
func (r *postRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return errors.WrapError(err, "Failed to delete post")
	}

	return nil
}

// GetPublished gets published posts
func (r *postRepository) GetPublished(limit, offset int) ([]*models.Post, error) {
	query := `SELECT id, title, content, author_id, is_published, created_at, updated_at 
			  FROM posts WHERE is_published = true 
			  ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get published posts")
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.IsPublished, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to scan post")
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Count returns the total number of posts
func (r *postRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, errors.WrapError(err, "Failed to count posts")
	}

	return count, nil
}

// CountByAuthorID returns the total number of posts by author
func (r *postRepository) CountByAuthorID(authorID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts WHERE author_id = $1`

	err := r.db.QueryRow(query, authorID).Scan(&count)
	if err != nil {
		return 0, errors.WrapError(err, "Failed to count posts by author")
	}

	return count, nil
}

// CountPublished returns the total number of published posts
func (r *postRepository) CountPublished() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts WHERE is_published = true`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, errors.WrapError(err, "Failed to count published posts")
	}

	return count, nil
}
