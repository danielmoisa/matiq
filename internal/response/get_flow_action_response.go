package response

import (
	"time"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type GetFlowActionResponse struct {
	UID         uuid.UUID              `json:"uid"`
	ResourceID  string                 `json:"resourceID,omitempty"`
	DisplayName string                 `json:"displayName"`
	ActionType  string                 `json:"actionType"`
	Template    map[string]interface{} `json:"template"`
	Transformer map[string]interface{} `json:"transformer"`
	TriggerMode string                 `json:"triggerMode"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"createdAt,omitempty"`
	CreatedBy   uuid.UUID              `json:"createdBy,omitempty"`
	UpdatedAt   time.Time              `json:"updatedAt,omitempty"`
	UpdatedBy   uuid.UUID              `json:"updatedBy,omitempty"`
}

func NewGetFlowActionResponse(flowAction *model.FlowAction) *GetFlowActionResponse {
	resp := &GetFlowActionResponse{
		UID:         flowAction.UID,
		ResourceID:  idconvertor.ConvertIntToString(flowAction.ResourceID),
		DisplayName: flowAction.Name,
		ActionType:  resourcelist.GetResourceIDMappedType(flowAction.ActionType),
		Template:    flowAction.ExportTemplateInMap(),
		Transformer: flowAction.ExportTransformerInMap(),
		TriggerMode: flowAction.TriggerMode,
		Config:      flowAction.ExportConfigInMap(),
		CreatedAt:   flowAction.CreatedAt,
		CreatedBy:   flowAction.CreatedBy,
		UpdatedAt:   flowAction.UpdatedAt,
		UpdatedBy:   flowAction.UpdatedBy,
	}
	return resp
}

func (resp *GetFlowActionResponse) ExportForFeedback() interface{} {
	return resp
}
