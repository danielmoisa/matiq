package response

import (
	"time"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type UpdateFlowActionResponse struct {
	UID         uuid.UUID              `json:"uid"`
	ResourceID  string                 `json:"resourceID,omitempty"`
	DisplayName string                 `json:"displayName"`
	ActionType  string                 `json:"actionType"`
	Template    map[string]interface{} `json:"template"`
	Transformer map[string]interface{} `json:"transformer"`
	TriggerMode string                 `json:"triggerMode"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"createdAt"`
	CreatedBy   string                 `json:"createdBy"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	UpdatedBy   string                 `json:"updatedBy"`
}

func NewUpdateFlowActionResponse(flowAction *model.FlowAction) *UpdateFlowActionResponse {
	resp := &UpdateFlowActionResponse{
		UID:         flowAction.UID,
		ResourceID:  idconvertor.ConvertIntToString(flowAction.ResourceID),
		DisplayName: flowAction.Name,
		ActionType:  resourcelist.GetResourceIDMappedType(flowAction.ActionType),
		Template:    flowAction.ExportTemplateInMap(),
		Transformer: flowAction.ExportTransformerInMap(),
		TriggerMode: flowAction.TriggerMode,
		Config:      flowAction.ExportConfigInMap(),
		CreatedAt:   flowAction.CreatedAt,
		CreatedBy:   flowAction.CreatedBy.String(),
		UpdatedAt:   flowAction.UpdatedAt,
		UpdatedBy:   flowAction.UpdatedBy.String(),
	}
	return resp
}

func (resp *UpdateFlowActionResponse) ExportForFeedback() interface{} {
	return resp
}
