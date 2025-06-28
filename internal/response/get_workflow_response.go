package response

import (
	"time"

	"github.com/danielmoisa/workflow-builder/internal/model"
	"github.com/danielmoisa/workflow-builder/internal/utils/idconvertor"
	"github.com/danielmoisa/workflow-builder/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type GetWorkflowResponse struct {
	ID                string                 `json:"workflowID"`
	UID               uuid.UUID              `json:"uid"`
	TeamID            string                 `json:"teamID"`
	Version           int                    `json:"version"`
	WorkflowID        string                 `json:"workflowID"`
	ResourceID        string                 `json:"resourceID,omitempty"`
	DisplayName       string                 `json:"displayName"`
	WorkflowType      string                 `json:"workflowType"`
	IsVirtualResource bool                   `json:"isVirtualResource"`
	Content           map[string]interface{} `json:"content"`
	Transformer       map[string]interface{} `json:"transformer"`
	TriggerMode       string                 `json:"triggerMode"`
	Config            map[string]interface{} `json:"config"`
	CreatedAt         time.Time              `json:"createdAt,omitempty"`
	CreatedBy         string                 `json:"createdBy,omitempty"`
	UpdatedAt         time.Time              `json:"updatedAt,omitempty"`
	UpdatedBy         string                 `json:"updatedBy,omitempty"`
}

func NewGetWorkflowResponse(workflow *model.Workflow) *GetWorkflowResponse {
	// workflowConfig := workflow.ExportConfig()
	resp := &GetWorkflowResponse{
		ID:           idconvertor.ConvertIntToString(workflow.ID),
		UID:          workflow.UID,
		TeamID:       idconvertor.ConvertIntToString(workflow.TeamID),
		WorkflowID:   idconvertor.ConvertIntToString(workflow.ID),
		Version:      workflow.Version,
		ResourceID:   idconvertor.ConvertIntToString(workflow.ResourceID),
		DisplayName:  workflow.Name,
		WorkflowType: resourcelist.GetResourceIDMappedType(workflow.Type),
		// IsVirtualResource: workflowConfig.IsVirtualResource,
		Content:     workflow.ExportTemplateInMap(),
		Transformer: workflow.ExportTransformerInMap(),
		TriggerMode: workflow.TriggerMode,
		Config:      workflow.ExportConfigInMap(),
		CreatedAt:   workflow.CreatedAt,
		CreatedBy:   idconvertor.ConvertIntToString(workflow.CreatedBy),
		UpdatedAt:   workflow.UpdatedAt,
		UpdatedBy:   idconvertor.ConvertIntToString(workflow.UpdatedBy),
	}
	return resp
}

func (resp *GetWorkflowResponse) ExportForFeedback() interface{} {
	return resp
}

// func (req *GetWorkflowResponse) AppendVirtualResourceToTemplate(value interface{}) {
// 	req.Content[model.ACTION_CONFIG_FIELD_VIRTUAL_RESOURCE] = value
// }
