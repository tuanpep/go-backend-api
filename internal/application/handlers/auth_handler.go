package handlers

import (
	"go-backend-api/internal/domain/entities"
	"go-backend-api/internal/pkg/auth"
	"go-backend-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService entities.UserService
	jwtManager  *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService entities.UserService, jwtManager *auth.JWTManager) *AuthHandler {
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
// @Param        request  body      entities.CreateUserRequest  true  "User registration data"
// @Success      201      {object}  response.Response{data=entities.User}
// @Failure      400      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req entities.CreateUserRequest
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
// @Param        request  body      entities.LoginRequest  true  "Login credentials"
// @Success      200      {object}  response.Response{data=entities.LoginResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entities.LoginRequest
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
