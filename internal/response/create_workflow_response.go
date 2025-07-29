package response

import (
	"time"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type CreateWorkflowResponse struct {
	FlowActionID      string                 `json:"flowActionID"`
	UID               uuid.UUID              `json:"uid"`
	ResourceID        string                 `json:"resourceID,omitempty"`
	DisplayName       string                 `json:"displayName"`
	ActionType        string                 `json:"flowActionType"`
	IsVirtualResource bool                   `json:"isVirtualResource"`
	Content           map[string]interface{} `json:"content"`
	Transformer       map[string]interface{} `json:"transformer"`
	TriggerMode       string                 `json:"triggerMode"`
	Config            map[string]interface{} `json:"config"`
	CreatedAt         time.Time              `json:"createdAt,omitempty"`
	CreatedBy         uuid.UUID              `json:"createdBy,omitempty"`
	UpdatedAt         time.Time              `json:"updatedAt,omitempty"`
	UpdatedBy         uuid.UUID              `json:"updatedBy,omitempty"`
}

func NewCreateWorkflowResponse(workflow *model.Workflow) *CreateWorkflowResponse {
	// flowActionConfig := workflow.ExportConfig()
	resp := &CreateWorkflowResponse{
		// FlowActionID: idconvertor.ConvertIntToString(workflow.UID),
		UID:         workflow.UID,
		ResourceID:  idconvertor.ConvertIntToString(workflow.ResourceID),
		DisplayName: workflow.Name,
		ActionType:  resourcelist.GetResourceIDMappedType(workflow.Type),
		// IsVirtualResource: flowActionConfig.IsVirtualResource,
		Content:     workflow.ExportTemplateInMap(),
		Transformer: workflow.ExportTransformerInMap(),
		TriggerMode: workflow.TriggerMode,
		Config:      workflow.ExportConfigInMap(),
		CreatedAt:   workflow.CreatedAt,
		CreatedBy:   workflow.CreatedBy,
		UpdatedAt:   workflow.UpdatedAt,
		UpdatedBy:   workflow.UpdatedBy,
	}
	return resp
}

func (resp *CreateWorkflowResponse) ExportForFeedback() interface{} {
	return resp
}
