package handlers

import (
	"go-backend-api/internal/models"
	"go-backend-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user requests
type UserHandler struct {
	userService models.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService models.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile gets the current user's profile
// @Summary      Get user profile
// @Description  Get the authenticated user's profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=models.User}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
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

	user, err := h.userService.GetUserByID(userUUID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// GetMe gets the current user's information
// @Summary      Get current user information
// @Description  Get the authenticated user's information (alias for /users/profile)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=models.User}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
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

	user, err := h.userService.GetUserByID(userUUID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateProfile updates the current user's profile
// @Summary      Update user profile
// @Description  Update the authenticated user's profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.UpdateUserRequest  true  "User update data"
// @Success      200      {object}  response.Response{data=models.User}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	user, err := h.userService.UpdateUser(userUUID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Profile updated successfully", user)
}

// DeleteProfile deletes the current user's account
// @Summary      Delete user account
// @Description  Delete the authenticated user's account
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/profile [delete]
func (h *UserHandler) DeleteProfile(c *gin.Context) {
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

	err := h.userService.DeleteUser(userUUID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}

// ActivateUser activates a user account
// @Summary      Activate user account
// @Description  Activate a user account by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id}/activate [put]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.userService.ActivateUser(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "User activated successfully", nil)
}

// DeactivateUser deactivates a user account
// @Summary      Deactivate user account
// @Description  Deactivate a user account by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id}/deactivate [put]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.userService.DeactivateUser(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "User deactivated successfully", nil)
}

// Logout logs out the current user
// @Summary      Logout user
// @Description  Logout the authenticated user by revoking refresh token
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
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

	// Extract token_id from claims
	claimsInterface, exists := c.Get("claims")
	if !exists {
		response.Unauthorized(c, "Token claims not found")
		return
	}

	claims, ok := claimsInterface.(*models.TokenClaims)
	if !ok {
		response.Unauthorized(c, "Invalid token claims")
		return
	}

	// Revoke refresh token
	err := h.userService.Logout(userUUID, claims.TokenID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Logged out successfully", nil)
}
