package handlers

import (
	"strconv"

	"go-backend-api/internal/domain/entities"
	"go-backend-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostHandler handles post requests
type PostHandler struct {
	postService entities.PostService
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService entities.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// Create creates a new post
// @Summary      Create a new post
// @Description  Create a new post (authenticated users only)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      entities.CreatePostRequest  true  "Post data"
// @Success      201      {object}  response.Response{data=entities.Post}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	var req entities.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	post, err := h.postService.CreatePost(userUUID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, post)
}

// GetAll gets all posts with pagination
// @Summary      Get all posts
// @Description  Get all posts with pagination support
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        per_page  query     int     false  "Items per page"  default(10)
// @Param        author_id query     string  false  "Filter by author ID"
// @Success      200       {object}  response.PaginatedResponse{data=[]entities.Post}
// @Failure      400       {object}  response.Response
// @Failure      401       {object}  response.Response
// @Failure      500       {object}  response.Response
// @Router       /posts [get]
func (h *PostHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	authorID := c.Query("author_id")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	var posts []*entities.Post
	var total int
	var err error

	if authorID != "" {
		authorUUID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			response.BadRequest(c, "Invalid author_id")
			return
		}
		posts, total, err = h.postService.GetPostsByAuthor(authorUUID, page, perPage)
	} else {
		posts, total, err = h.postService.GetPosts(page, perPage)
	}

	if err != nil {
		response.Error(c, err)
		return
	}

	totalPages := (total + perPage - 1) / perPage
	meta := response.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	response.Paginated(c, posts, meta)
}

// GetByID gets a post by ID
// @Summary      Get post by ID
// @Description  Get a specific post by its ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  response.Response{data=entities.Post}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, post)
}

// Update updates a post
// @Summary      Update a post
// @Description  Update a post (author only)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "Post ID"
// @Param        request  body      entities.UpdatePostRequest  true  "Post update data"
// @Success      200      {object}  response.Response{data=entities.Post}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      403      {object}  response.Response
// @Failure      404      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	var req entities.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	post, err := h.postService.UpdatePost(postID, userUUID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Post updated successfully", post)
}

// Delete deletes a post
// @Summary      Delete a post
// @Description  Delete a post (author only)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	err = h.postService.DeletePost(postID, userUUID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Post deleted successfully", nil)
}
