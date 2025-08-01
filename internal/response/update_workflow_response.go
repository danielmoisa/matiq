package response

import (
	"time"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type UpdateWorkflowResponse struct {
	WorkflowID        string                 `json:"workflowID"`
	UID               uuid.UUID              `json:"uid"`
	TeamID            string                 `json:"teamID,omitempty"`
	Version           int                    `json:"version"`
	ResourceID        string                 `json:"resourceID,omitempty"`
	DisplayName       string                 `json:"displayName"`
	WorkflowType      string                 `json:"workflowType"`
	IsVirtualResource bool                   `json:"isVirtualResource"`
	Template          map[string]interface{} `json:"template"`
	Transformer       map[string]interface{} `json:"transformer"`
	TriggerMode       string                 `json:"triggerMode"`
	Config            map[string]interface{} `json:"config"`
	CreatedAt         time.Time              `json:"createdAt"`
	CreatedBy         string                 `json:"createdBy"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	UpdatedBy         string                 `json:"updatedBy"`
}

func NewUpdateWorkflowResponse(workflow *model.Workflow) *UpdateWorkflowResponse {
	resp := &UpdateWorkflowResponse{
		WorkflowID:        idconvertor.ConvertIntToString(workflow.ResourceID),
		UID:               workflow.UID,
		Version:           model.WORKFLOW_EDIT_VERSION,
		ResourceID:        idconvertor.ConvertIntToString(workflow.ResourceID),
		DisplayName:       workflow.Name,
		WorkflowType:      resourcelist.GetResourceIDMappedType(workflow.Type),
		IsVirtualResource: resourcelist.IsVirtualResource(resourcelist.GetResourceIDMappedType(workflow.Type)),
		Template:          workflow.ExportTemplateInMap(),
		Transformer:       workflow.ExportTransformerInMap(),
		TriggerMode:       workflow.TriggerMode,
		Config:            workflow.ExportConfigInMap(),
		CreatedAt:         workflow.CreatedAt,
		CreatedBy:         workflow.CreatedBy.String(),
		UpdatedAt:         workflow.UpdatedAt,
		UpdatedBy:         workflow.UpdatedBy.String(),
	}
	return resp
}

func (resp *UpdateWorkflowResponse) ExportForFeedback() interface{} {
	return resp
}
