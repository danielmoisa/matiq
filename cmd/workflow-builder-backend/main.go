package main

import (
	"os"

	"github.com/danielmoisa/workflow-builder/internal/config"
	"github.com/danielmoisa/workflow-builder/internal/controller"
	"github.com/danielmoisa/workflow-builder/internal/driver/postgres"
	"github.com/danielmoisa/workflow-builder/internal/driver/redis"
	"github.com/danielmoisa/workflow-builder/internal/repository"
	"github.com/danielmoisa/workflow-builder/internal/router"
	"github.com/danielmoisa/workflow-builder/internal/utils/cache"
	"github.com/danielmoisa/workflow-builder/internal/utils/cors"
	"github.com/danielmoisa/workflow-builder/internal/utils/logger"
	"github.com/danielmoisa/workflow-builder/internal/utils/recovery"
	"github.com/gin-gonic/gin"

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
		logger.Errorw("Error in startup, repository init failed.")
	}
	return repository.NewRepository(postgresDriver, logger)
}

func initCache(globalConfig *config.Config, logger *zap.SugaredLogger) *cache.Cache {
	redisDriver, err := redis.NewRedisConnectionByGlobalConfig(globalConfig, logger)
	if err != nil {
		logger.Errorw("Error in startup, cache init failed.")
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

func initServer() (*Server, error) {
	globalConfig := config.GetInstance()
	engine := gin.New()
	sugaredLogger := logger.NewSugardLogger()

	// init validator
	// validator := tokenvalidator.NewRequestTokenValidator()

	// init driver
	repository := initRepo(globalConfig, sugaredLogger)
	cache := initCache(globalConfig, sugaredLogger)
	// drive := initDrive(globalConfig, sugaredLogger)

	// init attribute group
	// attrg, errInNewAttributeGroup := accesscontrol.NewRawAttributeGroup()
	// if errInNewAttributeGroup != nil {
	// 	return nil, errInNewAttributeGroup
	// }

	// init controller
	c := controller.NewControllerForBackend(repository, cache)
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
	server.router.RegisterRouters(server.engine)

	// run
	err := server.engine.Run(server.config.ServerHost + ":" + server.config.ServerPort)
	if err != nil {
		server.logger.Errorw("Error in startup", "err", err)
		os.Exit(2)
	}
}

func main() {
	server, err := initServer()

	if err != nil {

	}

	server.Start()
}
