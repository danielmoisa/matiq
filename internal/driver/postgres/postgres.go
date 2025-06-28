package postgres

import (
	"fmt"
	"time"

	"github.com/danielmoisa/workflow-builder/internal/config"
	"github.com/danielmoisa/workflow-builder/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const RETRY_TIMES = 6

type PostgresConfig struct {
	Addr     string `env:"WF_PG_ADDR" envDefault:"localhost"`
	Port     string `env:"WF_PG_PORT" envDefault:"5433"`
	User     string `env:"WF_PG_USER" envDefault:"wf_cloud"`
	Password string `env:"WF_PG_PASSWORD" envDefault:"wf2025"`
	Database string `env:"WF_PG_DATABASE" envDefault:"wf_cloud"`
}

func NewPostgresConnectionByGlobalConfig(config *config.Config, logger *zap.SugaredLogger) (*gorm.DB, error) {
	postgresConfig := &PostgresConfig{
		Addr:     config.GetPostgresAddr(),
		Port:     config.GetPostgresPort(),
		User:     config.GetPostgresUser(),
		Password: config.GetPostgresPassword(),
		Database: config.GetPostgresDatabase(),
	}
	return NewPostgresConnection(postgresConfig, logger)
}

func NewPostgresConnection(config *PostgresConfig, logger *zap.SugaredLogger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	retries := RETRY_TIMES

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host='%s' user='%s' password='%s' dbname='%s' port='%s'",
			config.Addr, config.User, config.Password, config.Database, config.Port),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	for err != nil {
		if logger != nil {
			logger.Errorw("Failed to connect to database, %d", retries)
		}
		if retries > 1 {
			retries--
			time.Sleep(10 * time.Second)
			db, err = gorm.Open(postgres.New(postgres.Config{
				DSN: fmt.Sprintf("host='%s' user='%s' password='%s' dbname='%s' port='%s'",
					config.Addr, config.User, config.Password, config.Database, config.Port),
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), &gorm.Config{})
			continue
		}
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw("error in connecting db ", "db", config, "err", err)
		return nil, err
	}

	// check db connection
	err = sqlDB.Ping()
	if err != nil {
		logger.Errorw("error in connecting db ", "db", config, "err", err)
		return nil, err
	}

	logger.Infow("connected with db", "db", config)

	// Run auto migrations
	err = runAutoMigrations(db, logger)
	if err != nil {
		logger.Errorw("error running auto migrations", "err", err)
		return nil, err
	}

	return db, err
}

// GetMigrationModels returns all models that should be migrated
func GetMigrationModels() []interface{} {
	return []interface{}{
		&model.Workflow{},
	}
}

// runAutoMigrations runs GORM auto migrations for all registered models
func runAutoMigrations(db *gorm.DB, logger *zap.SugaredLogger) error {
	models := GetMigrationModels()

	logger.Infow("Starting auto migrations", "models_count", len(models))

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			logger.Errorw("failed to migrate model", "model", fmt.Sprintf("%T", model), "err", err)
			return err
		}
		logger.Infow("successfully migrated model", "model", fmt.Sprintf("%T", model))
	}

	logger.Infow("Auto migrations completed successfully")
	return nil
}
