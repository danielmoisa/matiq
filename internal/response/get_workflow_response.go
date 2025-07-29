package response

import (
	"time"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type GetWorkflowResponse struct {
	UID               uuid.UUID              `json:"uid"`
	ResourceID        string                 `json:"resourceID,omitempty"`
	DisplayName       string                 `json:"displayName"`
	WorkflowType      string                 `json:"workflowType"`
	IsVirtualResource bool                   `json:"isVirtualResource"`
	Content           map[string]interface{} `json:"content"`
	Transformer       map[string]interface{} `json:"transformer"`
	TriggerMode       string                 `json:"triggerMode"`
	Template          map[string]interface{} `json:"template"`
	Config            map[string]interface{} `json:"config"`
	CreatedAt         time.Time              `json:"createdAt,omitempty"`
	CreatedBy         uuid.UUID              `json:"createdBy,omitempty"`
	UpdatedAt         time.Time              `json:"updatedAt,omitempty"`
	UpdatedBy         uuid.UUID              `json:"updatedBy,omitempty"`
}

func NewGetWorkflowResponse(workflow *model.Workflow) *GetWorkflowResponse {
	// workflowConfig := workflow.ExportConfig()
	resp := &GetWorkflowResponse{
		UID:          workflow.UID,
		ResourceID:   idconvertor.ConvertIntToString(workflow.ResourceID),
		DisplayName:  workflow.Name,
		WorkflowType: resourcelist.GetResourceIDMappedType(workflow.Type),
		// IsVirtualResource: workflowConfig.IsVirtualResource,
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

func (resp *GetWorkflowResponse) ExportForFeedback() interface{} {
	return resp
}

// func (req *GetWorkflowResponse) AppendVirtualResourceToTemplate(value interface{}) {
// 	req.Content[model.ACTION_CONFIG_FIELD_VIRTUAL_RESOURCE] = value
// }
