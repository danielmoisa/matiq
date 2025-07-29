package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env"
)

const DEPLOY_MODE_SELF_HOST = "self-host"
const DEPLOY_MODE_CLOUD = "cloud"
const DEPLOY_MODE_CLOUD_TEST = "cloud-test"
const DEPLOY_MODE_CLOUD_BETA = "cloud-beta"
const DEPLOY_MODE_CLOUD_PRODUCTION = "cloud-production"
const DRIVE_TYPE_AWS = "aws"
const DRIVE_TYPE_DO = "do"
const DRIVE_TYPE_MINIO = "minio"
const PROTOCOL_WEBSOCKET = "ws"
const PROTOCOL_WEBSOCKET_OVER_TLS = "wss"

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		var err error
		if instance == nil {
			instance, err = getConfig() // not thread safe
			if err != nil {
				panic(err)
			}
		}
	})
	return instance
}

type Config struct {
	// server config
	ServerHost         string `env:"MATIQ_SERVER_HOST"`
	ServerPort         string `env:"MATIQ_SERVER_PORT"`
	InternalServerPort string `env:"MATIQ_SERVER_INTERNAL_PORT"`
	ServerMode         string `env:"MATIQ_SERVER_MODE"`
	DeployMode         string `env:"MATIQ_DEPLOY_MODE"`
	SecretKey          string `env:"MATIQ_SECRET_KEY"`

	// websocket config
	WebsocketServerHost                       string `env:"MATIQ_WEBSOCKET_SERVER_HOST"`
	WebsocketServerPort                       string `env:"MATIQ_WEBSOCKET_SERVER_PORT"`
	WebsocketServerConnectionHost             string `env:"MATIQ_WEBSOCKET_CONNECTION_HOST"`
	WebsocketServerConnectionPort             string `env:"MATIQ_WEBSOCKET_CONNECTION_PORT"`
	WebsocketServerConnectionHostSouthAsia    string `env:"MATIQ_WEBSOCKET_CONNECTION_HOST_SOUTH_ASIA"`
	WebsocketServerConnectionPortSouthAsia    string `env:"MATIQ_WEBSOCKET_CONNECTION_PORT_SOUTH_ASIA"`
	WebsocketServerConnectionHostEastAsia     string `env:"MATIQ_WEBSOCKET_CONNECTION_HOST_EAST_ASIA"`
	WebsocketServerConnectionPortEastAsia     string `env:"MATIQ_WEBSOCKET_CONNECTION_PORT_EAST_ASIA"`
	WebsocketServerConnectionHostCenterEurope string `env:"MATIQ_WEBSOCKET_CONNECTION_HOST_CENTER_EUROPE"`
	WebsocketServerConnectionPortCenterEurope string `env:"MATIQ_WEBSOCKET_CONNECTION_PORT_CENTER_EUROPE"`
	WSSEnabled                                string `env:"MATIQ_WSS_ENABLED"`

	// key for idconvertor
	RandomKey string `env:"MATIQ_RANDOM_KEY"`
	// storage config
	PostgresAddr     string `env:"MATIQ_PG_ADDR"`
	PostgresPort     string `env:"MATIQ_PG_PORT"`
	PostgresUser     string `env:"MATIQ_PG_USER"`
	PostgresPassword string `env:"MATIQ_PG_PASSWORD"`
	PostgresDatabase string `env:"MATIQ_PG_DATABASE"`
	// cache config
	RedisAddr     string `env:"MATIQ_REDIS_ADDR"`
	RedisPort     string `env:"MATIQ_REDIS_PORT"`
	RedisPassword string `env:"MATIQ_REDIS_PASSWORD"`
	RedisDatabase int    `env:"MATIQ_REDIS_DATABASE"`
	// drive config
	DriveType             string `env:"MATIQ_DRIVE_TYPE"`
	DriveAccessKeyID      string `env:"MATIQ_DRIVE_ACCESS_KEY_ID"`
	DriveAccessKeySecret  string `env:"MATIQ_DRIVE_ACCESS_KEY_SECRET"`
	DriveRegion           string `env:"MATIQ_DRIVE_REGION"`
	DriveEndpoint         string `env:"MATIQ_DRIVE_ENDPOINT"`
	DriveSystemBucketName string `env:"MATIQ_DRIVE_SYSTEM_BUCKET_NAME"`
	DriveTeamBucketName   string `env:"MATIQ_DRIVE_TEAM_BUCKET_NAME"`
	DriveUploadTimeoutRaw string `env:"MATIQ_DRIVE_UPLOAD_TIMEOUT"`
	DriveUploadTimeout    time.Duration
	// supervisor API
	SupervisorInternalRestAPI string `env:"MATIQ_SUPERVISOR_INTERNAL_API"`

	// peripheral API
	PeripheralAPI string `env:"MATIQ_PERIPHERAL_API"`
	// resource manager API
	ResourceManagerRestAPI         string `env:"MATIQ_RESOURCE_MANAGER_API"`
	ResourceManagerInternalRestAPI string `env:"MATIQ_RESOURCE_MANAGER_INTERNAL_API"`
	// marketplace config
	MarketplaceInternalRestAPI string `env:"MATIQ_MARKETPLACE_INTERNAL_API"`
	// token for internal api
	ControlToken string `env:"MATIQ_CONTROL_TOKEN"`
	// google config
	GoogleSheetsClientID     string `env:"MATIQ_GS_CLIENT_ID"`
	GoogleSheetsClientSecret string `env:"MATIQ_GS_CLIENT_SECRET"`
	GoogleSheetsRedirectURI  string `env:"MATIQ_GS_REDIRECT_URI"`
	// toke for ip zone detector
	IPZoneDetectorToken string `env:"MATIQ_IP_ZONE_DETECTOR_TOKEN"`
	// drive config
	DriveRestAPI string `env:"MATIQ_DRIVE_API"`

	// Keycloak config
	KeycloakURL          string `env:"MATIQ_KEYCLOAK_URL"`
	KeycloakRealm        string `env:"MATIQ_KEYCLOAK_REALM"`
	KeycloakClientID     string `env:"MATIQ_KEYCLOAK_CLIENT_ID"`
	KeycloakClientSecret string `env:"MATIQ_KEYCLOAK_CLIENT_SECRET"`
	KeycloakAdminUser    string `env:"MATIQ_KEYCLOAK_ADMIN_USER"`
	KeycloakAdminPass    string `env:"MATIQ_KEYCLOAK_ADMIN_PASS"`
}

func getConfig() (*Config, error) {
	// Load .env file first
	loadEnvFile()

	// fetch
	cfg := &Config{}
	err := env.Parse(cfg)
	// process data
	var errInParseDuration error
	cfg.DriveUploadTimeout, errInParseDuration = time.ParseDuration(cfg.DriveUploadTimeoutRaw)
	if errInParseDuration != nil {
		return nil, errInParseDuration
	}
	// ok
	fmt.Printf("----------------\n")
	fmt.Printf("run by following config: %+v\n", cfg)
	fmt.Printf("parse config error info: %+v\n", err)

	return cfg, err
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env file is optional, so we don't error if it doesn't exist
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Only set if not already set in environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

func (c *Config) IsSelfHostMode() bool {
	return c.DeployMode == DEPLOY_MODE_SELF_HOST
}

func (c *Config) IsCloudMode() bool {
	if c.DeployMode == DEPLOY_MODE_CLOUD || c.DeployMode == DEPLOY_MODE_CLOUD_TEST || c.DeployMode == DEPLOY_MODE_CLOUD_BETA || c.DeployMode == DEPLOY_MODE_CLOUD_PRODUCTION {
		return true
	}
	return false
}

func (c *Config) IsCloudTestMode() bool {
	return c.DeployMode == DEPLOY_MODE_CLOUD_TEST
}

func (c *Config) IsCloudBetaMode() bool {
	return c.DeployMode == DEPLOY_MODE_CLOUD_BETA
}

func (c *Config) IsCloudProductionMode() bool {
	return c.DeployMode == DEPLOY_MODE_CLOUD_PRODUCTION
}

func (c *Config) GetWebScoketServerListenAddress() string {
	return c.WebsocketServerHost + ":" + c.WebsocketServerPort
}

func (c *Config) GetWebScoketServerConnectionAddress() string {
	return c.WebsocketServerConnectionHost + ":" + c.WebsocketServerConnectionPort
}

func (c *Config) GetWebsocketProtocol() string {
	if c.WSSEnabled == "true" {
		return PROTOCOL_WEBSOCKET_OVER_TLS
	}
	return PROTOCOL_WEBSOCKET
}

func (c *Config) GetRuntimeEnv() string {
	if c.IsCloudBetaMode() {
		return DEPLOY_MODE_CLOUD_BETA
	} else if c.IsCloudProductionMode() {
		return DEPLOY_MODE_CLOUD_PRODUCTION
	} else {
		return DEPLOY_MODE_CLOUD_TEST
	}
}

func (c *Config) GetSecretKey() string {
	return c.SecretKey
}

func (c *Config) GetRandomKey() string {
	return c.RandomKey
}

func (c *Config) GetPostgresAddr() string {
	return c.PostgresAddr
}

func (c *Config) GetPostgresPort() string {
	return c.PostgresPort
}

func (c *Config) GetPostgresUser() string {
	return c.PostgresUser
}

func (c *Config) GetPostgresPassword() string {
	return c.PostgresPassword
}

func (c *Config) GetPostgresDatabase() string {
	return c.PostgresDatabase
}

func (c *Config) GetRedisAddr() string {
	return c.RedisAddr
}

func (c *Config) GetRedisPort() string {
	return c.RedisPort
}

func (c *Config) GetRedisPassword() string {
	return c.RedisPassword
}

func (c *Config) GetRedisDatabase() int {
	return c.RedisDatabase
}

func (c *Config) GetDriveType() string {
	return c.DriveType
}

func (c *Config) IsAWSTypeDrive() bool {
	if c.DriveType == DRIVE_TYPE_AWS || c.DriveType == DRIVE_TYPE_DO {
		return true
	}
	return false
}

func (c *Config) IsMINIODrive() bool {
	return c.DriveType == DRIVE_TYPE_MINIO
}

func (c *Config) GetAWSS3Endpoint() string {
	return c.DriveEndpoint
}

func (c *Config) GetAWSS3AccessKeyID() string {
	return c.DriveAccessKeyID
}

func (c *Config) GetAWSS3AccessKeySecret() string {
	return c.DriveAccessKeySecret
}

func (c *Config) GetAWSS3Region() string {
	return c.DriveRegion
}

func (c *Config) GetAWSS3SystemBucketName() string {
	return c.DriveSystemBucketName
}

func (c *Config) GetAWSS3TeamBucketName() string {
	return c.DriveTeamBucketName
}

func (c *Config) GetAWSS3Timeout() time.Duration {
	return c.DriveUploadTimeout
}

func (c *Config) GetMINIOAccessKeyID() string {
	return c.DriveAccessKeyID
}

func (c *Config) GetMINIOAccessKeySecret() string {
	return c.DriveAccessKeySecret
}

func (c *Config) GetMINIOEndpoint() string {
	return c.DriveEndpoint
}

func (c *Config) GetMINIOSystemBucketName() string {
	return c.DriveSystemBucketName
}

func (c *Config) GetMINIOTeamBucketName() string {
	return c.DriveTeamBucketName
}

func (c *Config) GetMINIOTimeout() time.Duration {
	return c.DriveUploadTimeout
}

func (c *Config) GetControlToken() string {
	return c.ControlToken
}

func (c *Config) GetSupervisorInternalRestAPI() string {
	return c.SupervisorInternalRestAPI
}

func (c *Config) GetPeripheralAPI() string {
	return c.PeripheralAPI
}

func (c *Config) GetResourceManagerRestAPI() string {
	return c.ResourceManagerRestAPI
}

func (c *Config) GetResourceManagerInternalRestAPI() string {
	return c.ResourceManagerInternalRestAPI
}

func (c *Config) GetMarketplaceInternalRestAPI() string {
	return c.MarketplaceInternalRestAPI
}

func (c *Config) GetGoogleSheetsClientID() string {
	return c.GoogleSheetsClientID
}

func (c *Config) GetGoogleSheetsClientSecret() string {
	return c.GoogleSheetsClientSecret
}

func (c *Config) GetGoogleSheetsRedirectURI() string {
	return c.GoogleSheetsRedirectURI
}

func (c *Config) GetIPZoneDetectorToken() string {
	return c.IPZoneDetectorToken
}

func (c *Config) GetWebScoketServerConnectionAddressSouthAsia() string {
	return c.WebsocketServerConnectionHostSouthAsia + ":" + c.WebsocketServerConnectionPortSouthAsia
}

func (c *Config) GetWebScoketServerConnectionAddressEastAsia() string {
	return c.WebsocketServerConnectionHostEastAsia + ":" + c.WebsocketServerConnectionPortEastAsia
}

func (c *Config) GetWebScoketServerConnectionAddressCenterEurope() string {
	return c.WebsocketServerConnectionHostCenterEurope + ":" + c.WebsocketServerConnectionPortCenterEurope
}

func (c *Config) GetDriveAPIForSDK() string {
	return c.DriveRestAPI
}

// GetKeycloakConfig creates a Keycloak config from the main configuration
func (c *Config) GetKeycloakConfig() *KeycloakConfig {
	return &KeycloakConfig{
		URL:          c.KeycloakURL,
		Realm:        c.KeycloakRealm,
		ClientID:     c.KeycloakClientID,
		ClientSecret: c.KeycloakClientSecret,
		AdminUser:    c.KeycloakAdminUser,
		AdminPass:    c.KeycloakAdminPass,
	}
}

// KeycloakConfig represents Keycloak configuration
type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
	AdminUser    string
	AdminPass    string
}

// Validate checks if Keycloak configuration is valid
func (kc *KeycloakConfig) Validate() error {
	if kc.URL == "" {
		return fmt.Errorf("keycloak URL is required")
	}
	if kc.Realm == "" {
		return fmt.Errorf("keycloak realm is required")
	}
	if kc.ClientID == "" {
		return fmt.Errorf("keycloak client ID is required")
	}
	if kc.AdminUser == "" {
		return fmt.Errorf("keycloak admin user is required")
	}
	if kc.AdminPass == "" {
		return fmt.Errorf("keycloak admin password is required")
	}
	return nil
}

// IsEnabled checks if Keycloak is enabled (has basic required config)
func (c *Config) IsKeycloakEnabled() bool {
	return c.KeycloakURL != "" && c.KeycloakRealm != "" && c.KeycloakClientID != ""
}
