package request

import (
	"encoding/json"

	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
)

type UpdateFlowActionRequest struct {
	ID          string                 `json:"id"           validate:"required"`
	ActionType  string                 `json:"actionType"         validate:"required"`
	DisplayName string                 `json:"displayName"            validate:"required"`
	ResourceID  string                 `json:"resourceID,omitempty"`
	Template    map[string]interface{} `json:"template"                validate:"required"`
	Transformer map[string]interface{} `json:"transformer"            validate:"required"`
	TriggerMode string                 `json:"triggerMode"            validate:"oneof=manually automate"`
	Config      map[string]interface{} `json:"config"`
}

func NewUpdateFlowActionRequest() *UpdateFlowActionRequest {
	return &UpdateFlowActionRequest{}
}

func (req *UpdateFlowActionRequest) ExportTransformerInString() string {
	jsonByte, _ := json.Marshal(req.Transformer)
	return string(jsonByte)
}

func (req *UpdateFlowActionRequest) ExportFlowActionIDInInt() int {
	return idconvertor.ConvertStringToInt(req.ID)
}

func (req *UpdateFlowActionRequest) ExportResourceIDInInt() int {
	return idconvertor.ConvertStringToInt(req.ResourceID)
}

func (req *UpdateFlowActionRequest) ExportFlowActionTypeInInt() int {
	return resourcelist.GetResourceNameMappedID(req.ActionType)
}

func (req *UpdateFlowActionRequest) ExportTemplateInString() string {
	jsonByte, _ := json.Marshal(req.Template)
	return string(jsonByte)
}

func (req *UpdateFlowActionRequest) ExportConfigInString() string {
	jsonByte, _ := json.Marshal(req.Config)
	return string(jsonByte)
}

func (req *UpdateFlowActionRequest) AppendVirtualResourceToTemplate(value interface{}) {
	req.Template[ACTION_REQUEST_CONTENT_FIELD_VIRTUAL_RESOURCE] = value
}

func (req *UpdateFlowActionRequest) IsVirtualAction() bool {
	return resourcelist.IsVirtualResource(req.ActionType)
}

func (req *UpdateFlowActionRequest) IsLocalVirtualAction() bool {
	return resourcelist.IsLocalVirtualResource(req.ActionType)
}

func (req *UpdateFlowActionRequest) IsRemoteVirtualAction() bool {
	return resourcelist.IsRemoteVirtualResource(req.ActionType)
}

func (req *UpdateFlowActionRequest) NeedFetchResourceInfoFromSourceManager() bool {
	return resourcelist.NeedFetchResourceInfoFromSourceManager(req.ActionType)
}
