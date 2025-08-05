package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	FlowActionRepository *FlowActionRepository
}

func NewRepository(postgresDriver *gorm.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		FlowActionRepository: NewFlowActionRepository(logger, postgresDriver),
	}
}
