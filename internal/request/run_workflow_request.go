package request

import (
	"encoding/json"

	"github.com/danielmoisa/workflow-builder/internal/utils/idconvertor"
	"github.com/danielmoisa/workflow-builder/internal/utils/resourcelist"
)

// The run action HTTP request body like:
// ```json
//
//	{
//	    "resourceID": "ILAfx4p1C7dD",
//	    "actionType": "postgresql",
//	    "displayName": "postgresql1",
//	    "content": {
//	        "mode": "sql",
//	        "query": "select * from users where name like '%jame%';"
//	    },
//	    "context": {
//	        "input1.value": "jame"
//	    }
//	}
//
// ```

const (
	RUN_WORKFLOW_REQUEST_FIELD_CONTEXT = "context"
)

type RunWorkflowRequest struct {
	ResourceID   string                 `json:"resourceID,omitempty"`
	WorkflowType string                 `json:"workflowType" validate:"required"`
	DisplayName  string                 `json:"displayName" validate:"required"`
	Content      map[string]interface{} `json:"content" validate:"required"`
	Context      map[string]interface{} `json:"context" validate:"required"` // for action content raw param
}

func NewRunWorkflowRequest() *RunWorkflowRequest {
	return &RunWorkflowRequest{}
}

func (req *RunWorkflowRequest) ExportWorkflowTypeInInt() int {
	return resourcelist.GetResourceNameMappedID(req.WorkflowType)
}

func (req *RunWorkflowRequest) ExportResourceIDInInt() int {
	return idconvertor.ConvertStringToInt(req.ResourceID)
}

func (req *RunWorkflowRequest) ExportTemplateInString() string {
	jsonByte, _ := json.Marshal(req.Content)
	return string(jsonByte)
}

func (req *RunWorkflowRequest) ExportTemplateWithContextInString() string {
	content := req.Content
	content[RUN_WORKFLOW_REQUEST_FIELD_CONTEXT] = req.Content
	jsonByte, _ := json.Marshal(req.Content)
	return string(jsonByte)
}

func (req *RunWorkflowRequest) ExportContextInString() string {
	jsonByte, _ := json.Marshal(req.Context)
	return string(jsonByte)
}

func (req *RunWorkflowRequest) ExportContext() map[string]interface{} {
	return req.Context
}

func (req *RunWorkflowRequest) DoesContextAvaliable() bool {
	return len(req.Context) > 0
}

func (req *RunWorkflowRequest) IsVirtualAction() bool {
	return resourcelist.IsVirtualResource(req.WorkflowType)
}
