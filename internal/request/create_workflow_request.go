package request

import (
	"encoding/json"

	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
)

// The create action HTTP request body like:
// ```json
//
//	{
//	    "actionType": "postgresql",
//	    "displayName": "postgresql1",
//	    "resourceID": "ILAfx4p1C7cd",
//	    "content": {
//	        "mode": "sql",
//	        "query": ""
//	    },
//	    "isVirtualResource": true,
//	    "transformer": {
//	        "rawData": "",
//	        "enable": false
//	    },
//	    "triggerMode": "manually",
//	    "config": {
//	        "public": false,
//	        "advancedConfig": {
//	            "runtime": "none",
//	            "pages": [],
//	            "delayWhenLoaded": "",
//	            "displayLoadingPage": false,
//	            "isPeriodically": false,
//	            "periodInterval": ""
//	        }
//	    }
//	}
//
// ```
type CreateWorkflowRequest struct {
	WorkflowType      string                 `json:"workflowType" validate:"required"`
	DisplayName       string                 `json:"displayName" validate:"required"`
	ResourceID        string                 `json:"resourceID,omitempty"`
	IsVirtualResource bool                   `json:"isVirtualResource"`
	Template          map[string]interface{} `json:"template" validate:"required"`
	Transformer       map[string]interface{} `json:"transformer" validate:"required"`
	TriggerMode       string                 `json:"triggerMode" validate:"oneof=manually automate"`
	Config            map[string]interface{} `json:"config"`
}

func NewCreateWorkflowRequest() *CreateWorkflowRequest {
	return &CreateWorkflowRequest{}
}

func (req *CreateWorkflowRequest) ExportTransformerInString() string {
	jsonByte, _ := json.Marshal(req.Transformer)
	return string(jsonByte)
}

func (req *CreateWorkflowRequest) ExportFlowActionTypeInInt() int {
	return resourcelist.GetResourceNameMappedID(req.WorkflowType)
}

func (req *CreateWorkflowRequest) ExportResourceIDInInt() int {
	return idconvertor.ConvertStringToInt(req.ResourceID)
}

func (req *CreateWorkflowRequest) ExportTemplateInString() string {
	jsonByte, _ := json.Marshal(req.Template)
	return string(jsonByte)
}

func (req *CreateWorkflowRequest) ExportConfigInString() string {
	jsonByte, _ := json.Marshal(req.Config)
	return string(jsonByte)
}

func (req *CreateWorkflowRequest) AppendVirtualResourceToTemplate(value interface{}) {
	req.Template[ACTION_REQUEST_CONTENT_FIELD_VIRTUAL_RESOURCE] = value
}

func (req *CreateWorkflowRequest) IsVirtualAction() bool {
	return resourcelist.IsVirtualResource(req.WorkflowType)
}

func (req *CreateWorkflowRequest) IsLocalVirtualAction() bool {
	return resourcelist.IsLocalVirtualResource(req.WorkflowType)
}

func (req *CreateWorkflowRequest) IsRemoteVirtualAction() bool {
	return resourcelist.IsRemoteVirtualResource(req.WorkflowType)
}

func (req *CreateWorkflowRequest) NeedFetchResourceInfoFromSourceManager() bool {
	return resourcelist.NeedFetchResourceInfoFromSourceManager(req.WorkflowType)
}
