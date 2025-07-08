package keycloak

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/danielmoisa/auto-runner/internal/config"
)

// Config holds Keycloak configuration
type Config struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
	AdminUser    string
	AdminPass    string
}

// Client wraps the gocloak client with additional functionality
type Client struct {
	client      *gocloak.GoCloak
	config      *Config
	adminToken  *gocloak.JWT
	tokenExpiry time.Time
}

// UserInfo represents user information from Keycloak
type UserInfo struct {
	ID       string            `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Enabled  bool              `json:"enabled"`
	Roles    []string          `json:"roles"`
	Attrs    map[string]string `json:"attributes"`
}

// TokenInfo represents token validation result
type TokenInfo struct {
	Valid    bool                   `json:"valid"`
	UserID   string                 `json:"user_id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Enabled  bool                   `json:"enabled"`
	Roles    []string               `json:"roles"`
	Claims   map[string]interface{} `json:"claims"`
}

// NewClient creates a new Keycloak client
func NewClient(config *Config) *Client {
	client := gocloak.NewClient(config.URL)

	return &Client{
		client: client,
		config: config,
	}
}

// NewClientFromConfig creates a new Keycloak client from the main application config
func NewClientFromConfig(cfg *config.Config) (*Client, error) {
	keycloakConfig := cfg.GetKeycloakConfig()

	// Validate configuration
	if err := keycloakConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid keycloak configuration: %w", err)
	}

	// Convert to internal config format
	clientConfig := &Config{
		URL:          keycloakConfig.URL,
		Realm:        keycloakConfig.Realm,
		ClientID:     keycloakConfig.ClientID,
		ClientSecret: keycloakConfig.ClientSecret,
		AdminUser:    keycloakConfig.AdminUser,
		AdminPass:    keycloakConfig.AdminPass,
	}

	return NewClient(clientConfig), nil
}

// Connect establishes connection and gets admin token
func (c *Client) Connect(ctx context.Context) error {
	token, err := c.client.LoginAdmin(ctx, c.config.AdminUser, c.config.AdminPass, "master")
	if err != nil {
		return fmt.Errorf("failed to login admin: %w", err)
	}

	c.adminToken = token
	c.tokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	log.Printf("Connected to Keycloak realm: %s", c.config.Realm)
	return nil
}

// ensureAdminToken ensures we have a valid admin token
func (c *Client) ensureAdminToken(ctx context.Context) error {
	if c.adminToken == nil || time.Now().After(c.tokenExpiry.Add(-30*time.Second)) {
		return c.Connect(ctx)
	}
	return nil
}

// AuthenticateUser authenticates a user with username/password
func (c *Client) AuthenticateUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	token, err := c.client.Login(ctx, c.config.ClientID, c.config.ClientSecret, c.config.Realm, username, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return token, nil
}

// ValidateToken validates a JWT token and returns token info
func (c *Client) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	// First, introspect the token to check if it's active
	result, err := c.client.RetrospectToken(ctx, token, c.config.ClientID, c.config.ClientSecret, c.config.Realm)
	if err != nil {
		return &TokenInfo{Valid: false}, fmt.Errorf("token introspection failed: %w", err)
	}

	if !*result.Active {
		return &TokenInfo{Valid: false}, nil
	}

	// Get user info using the userinfo endpoint (this requires the access token)
	userInfo, err := c.client.GetUserInfo(ctx, token, c.config.Realm)
	if err != nil {
		log.Printf("Warning: failed to get user info from userinfo endpoint: %v", err)
		// Return basic validation without user details
		return &TokenInfo{
			Valid:    true,
			UserID:   "",
			Username: "",
			Email:    "",
			Enabled:  true,
			Roles:    []string{},
			Claims:   make(map[string]interface{}),
		}, nil
	}

	// Extract user information from userinfo response
	userID := getStringFromPointer(userInfo.Sub)
	username := getStringFromPointer(userInfo.PreferredUsername)
	if username == "" {
		username = getStringFromPointer(userInfo.Name)
	}
	email := getStringFromPointer(userInfo.Email)

	// Check if user is enabled (default to true if not specified)
	enabled := true
	if userInfo.EmailVerified != nil {
		enabled = *userInfo.EmailVerified
	}

	// Get user roles from Keycloak admin API using the user ID
	var roles []string
	if userID != "" {
		if err := c.ensureAdminToken(ctx); err == nil {
			userRoles, roleErr := c.client.GetRealmRolesByUserID(ctx, c.adminToken.AccessToken, c.config.Realm, userID)
			if roleErr == nil {
				for _, role := range userRoles {
					if role.Name != nil {
						roles = append(roles, *role.Name)
					}
				}
			} else {
				log.Printf("Warning: failed to get user roles: %v", roleErr)
			}
		} else {
			log.Printf("Warning: failed to ensure admin token: %v", err)
		}
	}

	return &TokenInfo{
		Valid:    true,
		UserID:   userID,
		Username: username,
		Email:    email,
		Enabled:  enabled,
		Roles:    roles,
		Claims:   convertUserInfoToMap(userInfo),
	}, nil
}

// GetUserInfo retrieves user information by user ID
func (c *Client) GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
	if err := c.ensureAdminToken(ctx); err != nil {
		return nil, err
	}

	user, err := c.client.GetUserByID(ctx, c.adminToken.AccessToken, c.config.Realm, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user roles
	roles, err := c.client.GetRealmRolesByUserID(ctx, c.adminToken.AccessToken, c.config.Realm, userID)
	if err != nil {
		log.Printf("Warning: failed to get user roles: %v", err)
	}

	var roleNames []string
	for _, role := range roles {
		if role.Name != nil {
			roleNames = append(roleNames, *role.Name)
		}
	}

	// Convert attributes
	attrs := make(map[string]string)
	if user.Attributes != nil {
		for key, values := range *user.Attributes {
			if len(values) > 0 {
				attrs[key] = values[0]
			}
		}
	}

	return &UserInfo{
		ID:       *user.ID,
		Username: *user.Username,
		Email:    getStringPtr(user.Email),
		Enabled:  getBoolPtr(user.Enabled),
		Roles:    roleNames,
		Attrs:    attrs,
	}, nil
}

// CreateUser creates a new user in Keycloak
func (c *Client) CreateUser(ctx context.Context, username, email, password string) (*UserInfo, error) {
	if err := c.ensureAdminToken(ctx); err != nil {
		return nil, err
	}

	user := gocloak.User{
		Username: &username,
		Email:    &email,
		Enabled:  gocloak.BoolP(true),
	}

	userID, err := c.client.CreateUser(ctx, c.adminToken.AccessToken, c.config.Realm, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Set password
	err = c.client.SetPassword(ctx, c.adminToken.AccessToken, c.config.Realm, userID, password, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	return c.GetUserInfo(ctx, userID)
}

// UpdateUser updates user information
func (c *Client) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	if err := c.ensureAdminToken(ctx); err != nil {
		return err
	}

	user, err := c.client.GetUserByID(ctx, c.adminToken.AccessToken, c.config.Realm, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Apply updates
	if email, ok := updates["email"].(string); ok {
		user.Email = &email
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		user.Enabled = &enabled
	}

	err = c.client.UpdateUser(ctx, c.adminToken.AccessToken, c.config.Realm, *user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user from Keycloak
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	if err := c.ensureAdminToken(ctx); err != nil {
		return err
	}

	err := c.client.DeleteUser(ctx, c.adminToken.AccessToken, c.config.Realm, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// AssignRole assigns a realm role to a user
func (c *Client) AssignRole(ctx context.Context, userID, roleName string) error {
	if err := c.ensureAdminToken(ctx); err != nil {
		return err
	}

	role, err := c.client.GetRealmRole(ctx, c.adminToken.AccessToken, c.config.Realm, roleName)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = c.client.AddRealmRoleToUser(ctx, c.adminToken.AccessToken, c.config.Realm, userID, []gocloak.Role{*role})
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RemoveRole removes a realm role from a user
func (c *Client) RemoveRole(ctx context.Context, userID, roleName string) error {
	if err := c.ensureAdminToken(ctx); err != nil {
		return err
	}

	role, err := c.client.GetRealmRole(ctx, c.adminToken.AccessToken, c.config.Realm, roleName)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = c.client.DeleteRealmRoleFromUser(ctx, c.adminToken.AccessToken, c.config.Realm, userID, []gocloak.Role{*role})
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	return nil
}

// RefreshToken refreshes an access token using a refresh token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	token, err := c.client.RefreshToken(ctx, refreshToken, c.config.ClientID, c.config.ClientSecret, c.config.Realm)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	return token, nil
}

// Logout logs out a user (invalidates tokens)
func (c *Client) Logout(ctx context.Context, refreshToken string) error {
	err := c.client.Logout(ctx, c.config.ClientID, c.config.ClientSecret, c.config.Realm, refreshToken)
	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	return nil
}

// HealthCheck verifies Keycloak connection
func (c *Client) HealthCheck(ctx context.Context) error {
	// Try to get server info as a health check
	_, err := c.client.GetServerInfo(ctx, c.adminToken.AccessToken)
	if err != nil {
		// If admin token is invalid, try to reconnect
		if err := c.Connect(ctx); err != nil {
			return fmt.Errorf("keycloak health check failed: %w", err)
		}
	}

	return nil
}

// GetUserGroups retrieves groups for a user
func (c *Client) GetUserGroups(ctx context.Context, userID string) ([]string, error) {
	if err := c.ensureAdminToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}

	groups, err := c.client.GetUserGroups(ctx, c.adminToken.AccessToken, c.config.Realm, userID, gocloak.GetGroupsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	var groupNames []string
	for _, group := range groups {
		if group.Name != nil {
			groupNames = append(groupNames, *group.Name)
		}
	}

	return groupNames, nil
}

// CheckUserInGroup checks if user is in a specific group
func (c *Client) CheckUserInGroup(ctx context.Context, userID, groupName string) (bool, error) {
	groups, err := c.GetUserGroups(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, group := range groups {
		if group == groupName {
			return true, nil
		}
	}

	return false, nil
}

// Helper functions
func getStringPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getBoolPtr(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

// Helper function to safely extract string from pointer
func getStringFromPointer(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// Helper function to convert UserInfo to map
func convertUserInfoToMap(userInfo *gocloak.UserInfo) map[string]interface{} {
	claims := make(map[string]interface{})

	if userInfo.Sub != nil {
		claims["sub"] = *userInfo.Sub
	}
	if userInfo.PreferredUsername != nil {
		claims["preferred_username"] = *userInfo.PreferredUsername
	}
	if userInfo.Email != nil {
		claims["email"] = *userInfo.Email
	}
	if userInfo.EmailVerified != nil {
		claims["email_verified"] = *userInfo.EmailVerified
	}
	if userInfo.Name != nil {
		claims["name"] = *userInfo.Name
	}
	if userInfo.GivenName != nil {
		claims["given_name"] = *userInfo.GivenName
	}
	if userInfo.FamilyName != nil {
		claims["family_name"] = *userInfo.FamilyName
	}

	return claims
}
