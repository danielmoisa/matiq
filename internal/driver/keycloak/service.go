package keycloak

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/danielmoisa/workflow-builder/internal/config"
)

// ServiceManager manages Keycloak service lifecycle
type ServiceManager struct {
	client *Client
	config *config.Config
}

// NewServiceManager creates a new Keycloak service manager
func NewServiceManager(cfg *config.Config) (*ServiceManager, error) {
	if !cfg.IsKeycloakEnabled() {
		return nil, fmt.Errorf("keycloak is not enabled in configuration")
	}

	client, err := NewClientFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create keycloak client: %w", err)
	}

	return &ServiceManager{
		client: client,
		config: cfg,
	}, nil
}

// Initialize connects to Keycloak and performs initial setup
func (sm *ServiceManager) Initialize(ctx context.Context) error {
	log.Printf("Initializing Keycloak connection to %s", sm.config.KeycloakURL)

	// Connect to Keycloak
	if err := sm.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to keycloak: %w", err)
	}

	// Perform health check
	if err := sm.client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("keycloak health check failed: %w", err)
	}

	log.Printf("Successfully connected to Keycloak realm: %s", sm.config.KeycloakRealm)
	return nil
}

// GetClient returns the Keycloak client
func (sm *ServiceManager) GetClient() *Client {
	return sm.client
}

// Shutdown gracefully shuts down the Keycloak service
func (sm *ServiceManager) Shutdown(ctx context.Context) error {
	log.Println("Shutting down Keycloak service...")
	// Add any cleanup logic here if needed
	return nil
}

// SetupDefaultRoles creates default roles if they don't exist
func (sm *ServiceManager) SetupDefaultRoles(ctx context.Context) error {
	defaultRoles := []string{"user", "admin", "workflow-creator", "workflow-executor"}

	for _, roleName := range defaultRoles {
		// Check if role exists, create if not
		_, err := sm.client.client.GetRealmRole(ctx, sm.client.adminToken.AccessToken, sm.config.KeycloakRealm, roleName)
		if err != nil {
			// Role doesn't exist, create it
			roleName := roleName
			role := gocloak.Role{
				Name:        &roleName,
				Description: gocloak.StringP(fmt.Sprintf("Default %s role", roleName)),
			}

			log.Printf("Creating default role: %s", roleName)
			_, err = sm.client.client.CreateRealmRole(ctx, sm.client.adminToken.AccessToken, sm.config.KeycloakRealm, role)
			if err != nil {
				log.Printf("Warning: Failed to create role %s: %v", roleName, err)
			} else {
				log.Printf("Successfully created role: %s", roleName)
			}
		}
	}

	return nil
}

// CreateDefaultUser creates a default admin user if it doesn't exist
func (sm *ServiceManager) CreateDefaultUser(ctx context.Context, username, email, password string, roles []string) error {
	// Check if user already exists
	params := gocloak.GetUsersParams{
		Username: &username,
	}
	users, err := sm.client.client.GetUsers(ctx, sm.client.adminToken.AccessToken, sm.config.KeycloakRealm, params)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	if len(users) > 0 {
		log.Printf("User %s already exists, skipping creation", username)
		return nil
	}

	// Create user
	log.Printf("Creating default user: %s", username)
	userInfo, err := sm.client.CreateUser(ctx, username, email, password)
	if err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	// Assign roles
	for _, role := range roles {
		err = sm.client.AssignRole(ctx, userInfo.ID, role)
		if err != nil {
			log.Printf("Warning: Failed to assign role %s to user %s: %v", role, username, err)
		} else {
			log.Printf("Assigned role %s to user %s", role, username)
		}
	}

	log.Printf("Successfully created default user: %s", username)
	return nil
}

// MonitorHealth continuously monitors Keycloak health
func (sm *ServiceManager) MonitorHealth(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sm.client.HealthCheck(ctx); err != nil {
				log.Printf("Keycloak health check failed: %v", err)
				// Attempt to reconnect
				if err := sm.client.Connect(ctx); err != nil {
					log.Printf("Failed to reconnect to Keycloak: %v", err)
				} else {
					log.Println("Successfully reconnected to Keycloak")
				}
			}
		}
	}
}

// GetConfiguration returns the current Keycloak configuration
func (sm *ServiceManager) GetConfiguration() *config.KeycloakConfig {
	return sm.config.GetKeycloakConfig()
}

// IsHealthy checks if Keycloak service is healthy
func (sm *ServiceManager) IsHealthy(ctx context.Context) bool {
	return sm.client.HealthCheck(ctx) == nil
}
