package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	WorkflowRepository *WorkflowRepository
}

func NewRepository(postgresDriver *gorm.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		WorkflowRepository: NewWorkflowRepository(logger, postgresDriver),
	}
}
