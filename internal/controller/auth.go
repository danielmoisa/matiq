package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielmoisa/matiq/internal/driver/keycloak"
	"github.com/gin-gonic/gin"
)

// LoginRequest represents login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *AuthenticationData `json:"data,omitempty"`
}

// AuthenticationData represents authentication data
type AuthenticationData struct {
	User         *UserInfo `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Enabled  bool     `json:"enabled"`
	Roles    []string `json:"roles"`
}

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ValidateTokenResponse represents token validation response
type ValidateTokenResponse struct {
	Success bool                `json:"success"`
	Valid   bool                `json:"valid"`
	Data    *keycloak.TokenInfo `json:"data,omitempty"`
	Message string              `json:"message,omitempty"`
}

// Login authenticates a user and returns JWT tokens
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (controller *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Authenticate user with Keycloak
	token, err := controller.KeycloakClient.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Validate token to get user info
	tokenInfo, err := controller.KeycloakClient.ValidateToken(c.Request.Context(), token.AccessToken)
	if err != nil || !tokenInfo.Valid {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Token validation failed",
		})
		return
	}

	// TokenInfo now includes all necessary user information, no need for additional GetUserInfo call
	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		Data: &AuthenticationData{
			User: &UserInfo{
				ID:       tokenInfo.UserID,
				Username: tokenInfo.Username,
				Email:    tokenInfo.Email,
				Enabled:  tokenInfo.Enabled,
				Roles:    tokenInfo.Roles,
			},
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    token.ExpiresIn,
			ExpiresAt:    expiresAt,
		},
	})
}

// Register creates a new user account
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (controller *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Create user in Keycloak
	userInfo, err := controller.KeycloakClient.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, LoginResponse{
				Success: false,
				Message: "User already exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "Failed to create user: " + err.Error(),
		})
		return
	}

	// Assign default user role
	err = controller.KeycloakClient.AssignRole(c.Request.Context(), userInfo.ID, "user")
	if err != nil {
		fmt.Printf("Warning: Could not assign default role: %v\n", err)
	}

	// Auto-login the newly created user
	token, err := controller.KeycloakClient.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusCreated, LoginResponse{
			Success: true,
			Message: "User created successfully. Please login.",
		})
		return
	}

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	c.JSON(http.StatusCreated, LoginResponse{
		Success: true,
		Message: "User created and logged in successfully",
		Data: &AuthenticationData{
			User: &UserInfo{
				ID:       userInfo.ID,
				Username: userInfo.Username,
				Email:    userInfo.Email,
				Enabled:  userInfo.Enabled,
				Roles:    userInfo.Roles,
			},
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    token.ExpiresIn,
			ExpiresAt:    expiresAt,
		},
	})
}

// RefreshToken refreshes an expired access token
// @Summary Refresh access token
// @Description Refresh an expired access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (controller *Controller) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Refresh token with Keycloak
	token, err := controller.KeycloakClient.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Invalid or expired refresh token",
		})
		return
	}

	// Validate new token to get user info
	tokenInfo, err := controller.KeycloakClient.ValidateToken(c.Request.Context(), token.AccessToken)
	if err != nil || !tokenInfo.Valid {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Token validation failed",
		})
		return
	}

	// TokenInfo now includes all necessary user information
	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: &AuthenticationData{
			User: &UserInfo{
				ID:       tokenInfo.UserID,
				Username: tokenInfo.Username,
				Email:    tokenInfo.Email,
				Enabled:  tokenInfo.Enabled,
				Roles:    tokenInfo.Roles,
			},
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    token.ExpiresIn,
			ExpiresAt:    expiresAt,
		},
	})
}

// Logout invalidates user tokens
// @Summary User logout
// @Description Logout user and invalidate tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Refresh token to invalidate"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/logout [post]
func (controller *Controller) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Logout from Keycloak
	err := controller.KeycloakClient.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Logout failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

// ValidateToken validates a JWT token
// @Summary Validate access token
// @Description Validate a JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} ValidateTokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/validate [get]
func (controller *Controller) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ValidateTokenResponse{
			Success: false,
			Valid:   false,
			Message: "Authorization header required",
		})
		return
	}

	// Validate token
	tokenInfo, err := controller.validateAuthToken(c.Request.Context(), authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ValidateTokenResponse{
			Success: false,
			Valid:   false,
			Message: "Token validation failed: " + err.Error(),
		})
		return
	}

	if !tokenInfo.Valid {
		c.JSON(http.StatusUnauthorized, ValidateTokenResponse{
			Success: false,
			Valid:   false,
			Message: "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, ValidateTokenResponse{
		Success: true,
		Valid:   true,
		Data:    tokenInfo,
		Message: "Token is valid",
	})
}

// GetProfile returns current user profile
// @Summary Get user profile
// @Description Get current authenticated user profile
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [get]
func (controller *Controller) GetProfile(c *gin.Context) {
	// Try to get token info from context first (set by auth middleware)
	if tokenInfo, exists := c.Get("token_info"); exists {
		tokenInfoStruct := tokenInfo.(*keycloak.TokenInfo)
		c.JSON(http.StatusOK, ProfileResponse{
			Success: true,
			Message: "Profile retrieved successfully",
			Data: &UserInfo{
				ID:       tokenInfoStruct.UserID,
				Username: tokenInfoStruct.Username,
				Email:    tokenInfoStruct.Email,
				Enabled:  tokenInfoStruct.Enabled,
				Roles:    tokenInfoStruct.Roles,
			},
		})
		return
	}

	// Fallback: get user info from token directly
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// Get detailed user information (only if token info not available in context)
	userInfo, err := controller.KeycloakClient.GetUserInfo(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get user profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProfileResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data: &UserInfo{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Email:    userInfo.Email,
			Enabled:  userInfo.Enabled,
			Roles:    userInfo.Roles,
		},
	})
}

// AuthMiddleware validates JWT tokens for protected routes
func (controller *Controller) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Validate token
		tokenInfo, err := controller.validateAuthToken(c.Request.Context(), authHeader)
		if err != nil || !tokenInfo.Valid {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", tokenInfo.UserID)
		c.Set("username", tokenInfo.Username)
		c.Set("user_email", tokenInfo.Email)
		c.Set("user_enabled", tokenInfo.Enabled)
		c.Set("user_roles", tokenInfo.Roles)
		c.Set("token_info", tokenInfo) // Store complete token info for efficiency

		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (controller *Controller) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "User roles not found",
			})
			c.Abort()
			return
		}

		roles := userRoles.([]string)
		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Message: fmt.Sprintf("Required role '%s' not found", requiredRole),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper method to validate authorization token
func (controller *Controller) validateAuthToken(ctx context.Context, authHeader string) (*keycloak.TokenInfo, error) {
	// Remove Bearer prefix if present
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	return controller.KeycloakClient.ValidateToken(ctx, token)
}

// GetUserIDFromKeycloakAuth gets the user ID from Keycloak authentication context
// This is compatible with the existing GetUserIDFromAuth method but returns a string
func (controller *Controller) GetUserIDFromKeycloakAuth(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, can not fetch user ID in it.")
		return "", fmt.Errorf("input missing user_id field")
	}
	userIDStr, ok := userID.(string)
	if !ok {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, user ID is not string type in it.")
		return "", fmt.Errorf("input user_id in wrong format")
	}
	return userIDStr, nil
}

// GetUsernameFromKeycloakAuth gets the username from Keycloak authentication context
func (controller *Controller) GetUsernameFromKeycloakAuth(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, can not fetch username in it.")
		return "", fmt.Errorf("input missing username field")
	}
	usernameStr, ok := username.(string)
	if !ok {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, username is not string type in it.")
		return "", fmt.Errorf("input username in wrong format")
	}
	return usernameStr, nil
}

// GetUserEmailFromKeycloakAuth gets the user email from Keycloak authentication context
func (controller *Controller) GetUserEmailFromKeycloakAuth(c *gin.Context) (string, error) {
	email, exists := c.Get("user_email")
	if !exists {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, can not fetch user email in it.")
		return "", fmt.Errorf("input missing user_email field")
	}
	emailStr, ok := email.(string)
	if !ok {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, user email is not string type in it.")
		return "", fmt.Errorf("input user_email in wrong format")
	}
	return emailStr, nil
}

// GetUserEnabledFromKeycloakAuth gets the user enabled status from Keycloak authentication context
func (controller *Controller) GetUserEnabledFromKeycloakAuth(c *gin.Context) (bool, error) {
	enabled, exists := c.Get("user_enabled")
	if !exists {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, can not fetch user enabled status in it.")
		return false, fmt.Errorf("input missing user_enabled field")
	}
	enabledBool, ok := enabled.(bool)
	if !ok {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, user enabled is not bool type in it.")
		return false, fmt.Errorf("input user_enabled in wrong format")
	}
	return enabledBool, nil
}

// GetTokenInfoFromKeycloakAuth gets the complete token info from Keycloak authentication context
func (controller *Controller) GetTokenInfoFromKeycloakAuth(c *gin.Context) (*keycloak.TokenInfo, error) {
	tokenInfo, exists := c.Get("token_info")
	if !exists {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, can not fetch token info in it.")
		return nil, fmt.Errorf("input missing token_info field")
	}
	tokenInfoStruct, ok := tokenInfo.(*keycloak.TokenInfo)
	if !ok {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "auth token invalid, token info is not correct type in it.")
		return nil, fmt.Errorf("input token_info in wrong format")
	}
	return tokenInfoStruct, nil
}

// ValidateKeycloakAuth validates the Keycloak token from the Authorization header
// This method can be used in place of existing token validation methods
func (controller *Controller) ValidateKeycloakAuth(c *gin.Context) (*keycloak.TokenInfo, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "authorization header missing.")
		return nil, fmt.Errorf("authorization header required")
	}

	tokenInfo, err := controller.validateAuthToken(c.Request.Context(), authHeader)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "token validation failed: "+err.Error())
		return nil, err
	}

	if !tokenInfo.Valid {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_TOKEN_FAILED, "invalid token.")
		return nil, fmt.Errorf("invalid token")
	}

	return tokenInfo, nil
}

// Common response types
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ProfileResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *UserInfo `json:"data,omitempty"`
}
