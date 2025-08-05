package request

import (
	"encoding/json"

	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
)

const (
	RUN_FLOW_ACTION_REQUEST_FIELD_CONTEXT = "context"
)

type RunFlowActionRequest struct {
	ResourceID  string                 `json:"resourceID,omitempty"`
	ActionType  string                 `json:"actionType" validate:"required"`
	DisplayName string                 `json:"displayName" validate:"required"`
	Content     map[string]interface{} `json:"content" validate:"required"`
	Context     map[string]interface{} `json:"context" validate:"required"`
}

func NewRunFlowActionRequest() *RunFlowActionRequest {
	return &RunFlowActionRequest{}
}

func (req *RunFlowActionRequest) ExportWorkflowTypeInInt() int {
	return resourcelist.GetResourceNameMappedID(req.ActionType)
}

func (req *RunFlowActionRequest) ExportResourceIDInInt() int {
	return idconvertor.ConvertStringToInt(req.ResourceID)
}

func (req *RunFlowActionRequest) ExportTemplateInString() string {
	jsonByte, _ := json.Marshal(req.Content)
	return string(jsonByte)
}

func (req *RunFlowActionRequest) ExportTemplateWithContextInString() string {
	content := req.Content
	content[RUN_FLOW_ACTION_REQUEST_FIELD_CONTEXT] = req.Content
	jsonByte, _ := json.Marshal(req.Content)
	return string(jsonByte)
}

func (req *RunFlowActionRequest) ExportContextInString() string {
	jsonByte, _ := json.Marshal(req.Context)
	return string(jsonByte)
}

func (req *RunFlowActionRequest) ExportContext() map[string]interface{} {
	return req.Context
}

func (req *RunFlowActionRequest) DoesContextAvaliable() bool {
	return len(req.Context) > 0
}

func (req *RunFlowActionRequest) IsVirtualAction() bool {
	return resourcelist.IsVirtualResource(req.ActionType)
}
