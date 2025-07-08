package repository

import (
	"errors"
	"fmt"

	"github.com/danielmoisa/auto-runner/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WorkflowRepository struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewWorkflowRepository(logger *zap.SugaredLogger, db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{
		logger: logger,
		db:     db,
	}
}

func (impl *WorkflowRepository) Create(action *model.Workflow) (int, error) {
	if err := impl.db.Create(action).Error; err != nil {
		return 0, err
	}
	return action.ID, nil
}

func (impl *WorkflowRepository) Delete(teamID int, workflowID int) error {
	if err := impl.db.Where("id = ? AND team_id = ?", workflowID, teamID).Delete(&model.Workflow{}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *WorkflowRepository) UpdateWholeFlowAction(action *model.Workflow) error {
	if err := impl.db.Model(action).Where("id = ?", action.ID).UpdateColumns(action).Error; err != nil {
		return err
	}
	return nil
}

func (impl *WorkflowRepository) RetrieveWorkflowByID(workflowID uuid.UUID) (*model.Workflow, error) {
	var workflow model.Workflow
	if err := impl.db.Where("uid = ?", workflowID).First(&workflow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workflow not found for workflowID %d", workflowID)
		}
		return nil, err
	}
	return &workflow, nil
}

func (impl *WorkflowRepository) RetrieveAll(teamID int, workflowID int, version int) ([]*model.Workflow, error) {
	var actions []*model.Workflow
	if err := impl.db.Where("team_id = ? AND workflow_id = ? AND version = ?", teamID, workflowID, version).Find(&actions).Error; err != nil {
		return nil, err
	}
	return actions, nil
}

func (impl *WorkflowRepository) RetrieveByType(teamID int, workflowID int, version int, actionType int) ([]*model.Workflow, error) {
	var actions []*model.Workflow
	if err := impl.db.Where("team_id = ? AND workflow_id = ? AND version = ? AND type = ? ", teamID, workflowID, version, actionType).Find(&actions).Error; err != nil {
		return nil, err
	}
	return actions, nil
}

func (impl *WorkflowRepository) RetrieveByID(teamID int, workflowID int) (*model.Workflow, error) {
	var action *model.Workflow
	if err := impl.db.Where("id = ? AND team_id = ?", workflowID, teamID).First(&action).Error; err != nil {
		return &model.Workflow{}, err
	}
	return action, nil
}

func (impl *WorkflowRepository) RetrieveFlowActionsByTeamIDWorkflowIDAndVersion(teamID int, workflowID int, version int) ([]*model.Workflow, error) {
	var actions []*model.Workflow
	if err := impl.db.Where("team_id = ? AND workflow_id = ? AND version = ?", teamID, workflowID, version).Find(&actions).Error; err != nil {
		return nil, err
	}
	return actions, nil
}

func (impl *WorkflowRepository) RetrieveFlowActionsByTeamIDWorkflowIDVersionAndType(teamID int, workflowID int, version int, actionType int) ([]*model.Workflow, error) {
	var actions []*model.Workflow
	if err := impl.db.Where("team_id = ? AND workflow_id = ? AND version = ? AND type = ?", teamID, workflowID, version, actionType).Find(&actions).Error; err != nil {
		return nil, err
	}
	return actions, nil
}

func (impl *WorkflowRepository) RetrieveFlowActionByTeamIDFlowActionID(teamID int, workflowID int) (*model.Workflow, error) {
	var action *model.Workflow
	if err := impl.db.Where("team_id = ? AND id = ?", teamID, workflowID).First(&action).Error; err != nil {
		return nil, err
	}
	return action, nil
}

func (impl *WorkflowRepository) DeleteFlowActionsByWorkflow(teamID int, workflowID int) error {
	if err := impl.db.Where("team_id = ? AND workflow_id = ?", teamID, workflowID).Delete(&model.Workflow{}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *WorkflowRepository) DeleteFlowActionByTeamIDAndFlowActionID(teamID int, workflowID int) error {
	if err := impl.db.Where("team_id = ? AND id = ?", teamID, workflowID).Delete(&model.Workflow{}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *WorkflowRepository) CountFlowActionByTeamID(teamID int) (int, error) {
	var count int64
	if err := impl.db.Model(&model.Workflow{}).Where("team_id = ?", teamID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (impl *WorkflowRepository) DeleteAllFlowActionsByTeamIDWorkflowIDAndVersion(teamID int, workflowID int, targetVersion int) error {
	if err := impl.db.Where("team_id = ? AND workflow_id = ? AND version = ?", teamID, workflowID, targetVersion).Delete(&model.Workflow{}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *WorkflowRepository) RetrieveAllWorkflowByUserID(userID string) ([]*model.Workflow, error) {
	var workflows []*model.Workflow
	if err := impl.db.Where("created_by = ?", userID).Find(&workflows).Error; err != nil {
		return nil, err
	}
	return workflows, nil
}
