package main

import (
	"os"

	"github.com/danielmoisa/workflow-builder/docs"
	_ "github.com/danielmoisa/workflow-builder/docs"
	"github.com/danielmoisa/workflow-builder/internal/config"
	"github.com/danielmoisa/workflow-builder/internal/controller"
	"github.com/danielmoisa/workflow-builder/internal/driver/keycloak"
	"github.com/danielmoisa/workflow-builder/internal/driver/postgres"
	"github.com/danielmoisa/workflow-builder/internal/driver/redis"
	"github.com/danielmoisa/workflow-builder/internal/repository"
	"github.com/danielmoisa/workflow-builder/internal/router"
	"github.com/danielmoisa/workflow-builder/internal/utils/cache"
	"github.com/danielmoisa/workflow-builder/internal/utils/cors"
	"github.com/danielmoisa/workflow-builder/internal/utils/logger"
	"github.com/danielmoisa/workflow-builder/internal/utils/recovery"
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

// func initDrive(globalConfig *config.Config, logger *zap.SugaredLogger) *drive.Drive {
// 	if globalConfig.IsAWSTypeDrive() {
// 		teamAWSConfig := awss3.NewTeamAwsConfigByGlobalConfig(globalConfig)
// 		teamDriveS3Instance := awss3.NewS3Drive(teamAWSConfig)
// 		return drive.NewDrive(teamDriveS3Instance, logger)
// 	}
// 	// failed
// 	logger.Errorw("Error in startup, drive init failed.")
// 	return nil
// }

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

	// init validator
	// validator := tokenvalidator.NewRequestTokenValidator()

	// init driver
	repository := initRepo(globalConfig, sugaredLogger)
	cache := initCache(globalConfig, sugaredLogger)
	keycloakClient := initKeycloak(globalConfig, sugaredLogger)

	// drive := initDrive(globalConfig, sugaredLogger)

	// init attribute group
	// attrg, errInNewAttributeGroup := accesscontrol.NewRawAttributeGroup()
	// if errInNewAttributeGroup != nil {
	// 	return nil, errInNewAttributeGroup
	// }

	// init controller
	c := controller.NewControllerForBackend(repository, cache, keycloakClient)
	router := router.NewRouter(c)
	server := NewServer(globalConfig, engine, router, sugaredLogger)
	return server, nil

}

func (server *Server) Start() {
	server.logger.Infow("Starting workflow-builder-backend...")

	// init
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

// @title			Go Flow Automation API
// @version		1.0
// @description	A auth and authorization application built with Go
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT
// @host			localhost:8001
// @BasePath		/api/v1
// @schemes		http https
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
