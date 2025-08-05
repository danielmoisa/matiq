package repository

import (
	"errors"
	"fmt"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FlowActionRepository struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewFlowActionRepository(logger *zap.SugaredLogger, db *gorm.DB) *FlowActionRepository {
	return &FlowActionRepository{
		logger: logger,
		db:     db,
	}
}

func (impl *FlowActionRepository) Create(action *model.FlowAction) (uuid.UUID, error) {
	if err := impl.db.Create(action).Error; err != nil {
		return uuid.Nil, err
	}
	return action.UID, nil
}

func (impl *FlowActionRepository) Delete(teamID int, flowActionID int) error {
	if err := impl.db.Where("id = ? AND team_id = ?", flowActionID, teamID).Delete(&model.FlowAction{}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *FlowActionRepository) Update(flowAction *model.FlowAction) error {
	if err := impl.db.Model(flowAction).Where("uid = ?", flowAction.UID).UpdateColumns(flowAction).Error; err != nil {
		return err
	}
	return nil
}

func (impl *FlowActionRepository) RetrieveFlowActionByID(flowActionID uuid.UUID) (*model.FlowAction, error) {
	var flowAction model.FlowAction
	if err := impl.db.Where("uid = ?", flowActionID).First(&flowAction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("flowAction not found for flowActionID %d", flowActionID)
		}
		return nil, err
	}
	return &flowAction, nil
}

func (impl *FlowActionRepository) RetrieveAllFlowActionsByUserID(userID string) ([]*model.FlowAction, error) {
	var flowActions []*model.FlowAction
	if err := impl.db.Where("created_by = ?", userID).Find(&flowActions).Error; err != nil {
		return nil, err
	}
	return flowActions, nil
}
