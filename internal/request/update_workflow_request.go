package request

import (
	"encoding/json"

	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
)

// The update action HTTP request body like:
// ```json
//
//	{
//	    "actionID": "ILAex4p1C7rD",
//	    "uid": "781f0ed4-62eb-4615-bd41-80bf2af8ceb4",
//	    "teamID": "ILAfx4p1C7bN",
//	    "resourceID": "ILAfx4p1C7cc",
//	    "displayName": "postgresql1",
//	    "actionType": "postgresql",
//	    "isVirtualResource": false,
//	    "content": {
//	        "mode": "sql",
//	        "query": "select * from data;"
//	    },
//	    "transformer": {
//	        "enable": false,
//	        "rawData": ""
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
//	    },
//	    "createdAt": "2023-08-25T10:18:21.914894943Z",
//	    "createdBy": "ILAfx4p1C7dX",
//	    "updatedAt": "2023-08-25T10:18:21.91489513Z",
//	    "updatedBy": "ILAfx4p1C7dX"
//	}
//
// ```
type UpdateWorkflowRequest struct {
	WorkflowID        string                 `json:"workflowID"           validate:"required"`
	WorkflowType      string                 `json:"workflowType"         validate:"required"`
	DisplayName       string                 `json:"displayName"            validate:"required"`
	ResourceID        string                 `json:"resourceID,omitempty"`
	IsVirtualResource bool                   `json:"isVirtualResource"`
	Content           map[string]interface{} `json:"content"                validate:"required"`
	Transformer       map[string]interface{} `json:"transformer"            validate:"required"`
	TriggerMode       string                 `json:"triggerMode"            validate:"oneof=manually automate"`
	Config            map[string]interface{} `json:"config"`
}

func NewUpdateWorkflowRequest() *UpdateWorkflowRequest {
	return &UpdateWorkflowRequest{}
}

func (req *UpdateWorkflowRequest) ExportTransformerInString() string {
	jsonByte, _ := json.Marshal(req.Transformer)
	return string(jsonByte)
}

func (req *UpdateWorkflowRequest) ExportFlowActionIDInInt() int {
	return idconvertor.ConvertStringToInt(req.WorkflowID)
}

func (req *UpdateWorkflowRequest) ExportResourceIDInInt() int {
	return idconvertor.ConvertStringToInt(req.ResourceID)
}

func (req *UpdateWorkflowRequest) ExportWorkflowTypeInInt() int {
	return resourcelist.GetResourceNameMappedID(req.WorkflowType)
}

func (req *UpdateWorkflowRequest) ExportTemplateInString() string {
	jsonByte, _ := json.Marshal(req.Content)
	return string(jsonByte)
}

func (req *UpdateWorkflowRequest) ExportConfigInString() string {
	jsonByte, _ := json.Marshal(req.Config)
	return string(jsonByte)
}

func (req *UpdateWorkflowRequest) AppendVirtualResourceToTemplate(value interface{}) {
	req.Content[ACTION_REQUEST_CONTENT_FIELD_VIRTUAL_RESOURCE] = value
}

func (req *UpdateWorkflowRequest) IsVirtualAction() bool {
	return resourcelist.IsVirtualResource(req.WorkflowType)
}

func (req *UpdateWorkflowRequest) IsLocalVirtualAction() bool {
	return resourcelist.IsLocalVirtualResource(req.WorkflowType)
}

func (req *UpdateWorkflowRequest) IsRemoteVirtualAction() bool {
	return resourcelist.IsRemoteVirtualResource(req.WorkflowType)
}

func (req *UpdateWorkflowRequest) NeedFetchResourceInfoFromSourceManager() bool {
	return resourcelist.NeedFetchResourceInfoFromSourceManager(req.WorkflowType)
}
