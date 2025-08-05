package main

import (
	"os"

	"github.com/danielmoisa/matiq/docs"
	_ "github.com/danielmoisa/matiq/docs"
	"github.com/danielmoisa/matiq/internal/config"
	"github.com/danielmoisa/matiq/internal/controller"
	"github.com/danielmoisa/matiq/internal/driver/keycloak"
	"github.com/danielmoisa/matiq/internal/driver/postgres"
	"github.com/danielmoisa/matiq/internal/driver/redis"
	"github.com/danielmoisa/matiq/internal/repository"
	"github.com/danielmoisa/matiq/internal/router"
	"github.com/danielmoisa/matiq/internal/utils/cache"
	"github.com/danielmoisa/matiq/internal/utils/cors"
	"github.com/danielmoisa/matiq/internal/utils/logger"
	"github.com/danielmoisa/matiq/internal/utils/recovery"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	router *router.Router
	logger *zap.SugaredLogger
	config *config.Config
}

func NewServer(config *config.Config, engine *gin.Engine, router *router.Router, logger *zap.SugaredLogger) *Server {
	return &Server{
		engine: engine,
		config: config,
		router: router,
		logger: logger,
	}
}

func initRepo(globalConfig *config.Config, logger *zap.SugaredLogger) *repository.Repository {
	postgresDriver, err := postgres.NewPostgresConnectionByGlobalConfig(globalConfig, logger)
	if err != nil {
		logger.Errorw("Error in startup, repository init failed.", "error", err)
		return nil
	}
	return repository.NewRepository(postgresDriver, logger)
}

func initCache(globalConfig *config.Config, logger *zap.SugaredLogger) *cache.Cache {
	redisDriver, err := redis.NewRedisConnectionByGlobalConfig(globalConfig, logger)
	if err != nil {
		logger.Errorw("Error in startup, cache init failed.", "error", err)
		return nil
	}
	return cache.NewCache(redisDriver, logger)
}

func initKeycloak(globalConfig *config.Config, logger *zap.SugaredLogger) *keycloak.Client {
	keycloakClient, err := keycloak.NewClientFromConfig(globalConfig)
	if err != nil {
		logger.Errorw("Error in startup, keycloak init failed.", "error", err)
		return nil
	}
	return keycloakClient
}

func initServer() (*Server, error) {
	globalConfig := config.GetInstance()
	engine := gin.New()
	sugaredLogger := logger.NewSugardLogger()

	initRepo := initRepo(globalConfig, sugaredLogger)
	initCache := initCache(globalConfig, sugaredLogger)
	keycloakClient := initKeycloak(globalConfig, sugaredLogger)

	c := controller.NewControllerForBackend(initRepo, initCache, keycloakClient)
	appRouter := router.NewRouter(c)
	server := NewServer(globalConfig, engine, appRouter, sugaredLogger)
	return server, nil

}

func (server *Server) Start() {
	server.logger.Infow("Starting matiq server...")

	gin.SetMode(server.config.ServerMode)

	// init cors
	server.engine.Use(gin.CustomRecovery(recovery.CorsHandleRecovery))
	server.engine.Use(cors.Cors())

	// init swagger
	docs.SwaggerInfo.Host = server.config.ServerHost + ":" + server.config.ServerPort
	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	server.router.RegisterRouters(server.engine)

	// run
	err := server.engine.Run(server.config.ServerHost + ":" + server.config.ServerPort)
	if err != nil {
		server.logger.Errorw("Error in startup", "err", err)
		os.Exit(2)
	}
}

// @title Matiq Automation API
// @version 1.0
// @description API for managing automation workflows and tasks.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8001
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	server, err := initServer()

	if err != nil {
		logger.NewSugardLogger().Errorw("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	server.Start()
}
