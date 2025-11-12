package handlers

import (
	"go-backend-api/internal/models"
	"go-backend-api/internal/pkg/auth"
	"go-backend-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService models.UserService
	jwtManager  *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService models.UserService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Register a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateUserRequest  true  "User registration data"
// @Success      201      {object}  response.Response{data=models.User}
// @Failure      400      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, user)
}

// Login handles user login
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login credentials"
// @Success      200      {object}  response.Response{data=models.LoginResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	loginResp, err := h.userService.AuthenticateUser(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, loginResp)
}

// Refresh handles refresh token requests
// @Summary      Refresh access token
// @Description  Refresh access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RefreshTokenRequest  true  "Refresh token"
// @Success      200      {object}  response.Response{data=models.LoginResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	loginResp, err := h.userService.RefreshToken(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, loginResp)
}
