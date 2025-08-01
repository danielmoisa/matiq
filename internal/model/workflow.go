package model

import (
	"encoding/json"
	"time"

	"github.com/danielmoisa/matiq/internal/request"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

const (
	WORKFLOW_EDIT_VERSION = 0
)

type Workflow struct {
	UID         uuid.UUID              `gorm:"column:uid;type:uuid;not null"`
	ResourceID  int                    `gorm:"column:resource_id;type:bigint;not null"`
	Name        string                 `gorm:"column:name;type:varchar;size:255;not null"`
	Type        int                    `gorm:"column:type;type:smallint;not null"`
	TriggerMode string                 `gorm:"column:trigger_mode;type:varchar;size:16;not null"`
	Transformer string                 `gorm:"column:transformer;type:jsonb"`
	Template    string                 `gorm:"column:template;type:jsonb"`
	RawTemplate string                 `gorm:"-" sql:"-"`
	Context     map[string]interface{} `gorm:"-" sql:"-"`
	Config      string                 `gorm:"column:config;type:jsonb"`
	CreatedAt   time.Time              `gorm:"column:created_at;type:timestamp;not null"`
	CreatedBy   uuid.UUID              `gorm:"column:created_by;type:uuid;not null"`
	UpdatedAt   time.Time              `gorm:"column:updated_at;type:timestamp;not null"`
	UpdatedBy   uuid.UUID              `gorm:"column:updated_by;type:uuid;not null"`
}

func NewWorkflow() *Workflow {
	return &Workflow{}
}

func NewWorkflowByCreateRequest(userID uuid.UUID, req *request.CreateWorkflowRequest) (*Workflow, error) {
	action := &Workflow{
		ResourceID:  idconvertor.ConvertStringToInt(req.ResourceID),
		Name:        req.DisplayName,
		Type:        resourcelist.GetResourceNameMappedID(req.WorkflowType),
		TriggerMode: req.TriggerMode,
		Transformer: req.ExportTransformerInString(),
		Template:    req.ExportTemplateInString(),
		Config:      req.ExportConfigInString(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	action.InitUID()
	action.InitCreatedAt()
	action.InitUpdatedAt()
	return action, nil
}

func NewWorkflowByUpdateRequest(userID uuid.UUID, req *request.UpdateWorkflowRequest) (*Workflow, error) {
	action := &Workflow{
		ResourceID:  idconvertor.ConvertStringToInt(req.ResourceID),
		Name:        req.DisplayName,
		Type:        resourcelist.GetResourceNameMappedID(req.WorkflowType),
		TriggerMode: req.TriggerMode,
		Transformer: req.ExportTransformerInString(),
		Template:    req.ExportTemplateInString(),
		Config:      req.ExportConfigInString(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	action.InitUID()
	action.InitCreatedAt()
	action.InitUpdatedAt()
	return action, nil
}

func NewWorkflowByRunRequest(userID uuid.UUID, req *request.RunWorkflowRequest) *Workflow {
	action := &Workflow{
		ResourceID:  idconvertor.ConvertStringToInt(req.ResourceID),
		Name:        req.DisplayName,
		Type:        resourcelist.GetResourceNameMappedID(req.WorkflowType),
		Template:    req.ExportTemplateInString(),
		RawTemplate: req.ExportTemplateWithContextInString(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	action.InitUID()
	action.InitCreatedAt()
	action.InitUpdatedAt()
	return action
}

func (action *Workflow) InitUID() {
	action.UID = uuid.New()
}

func (action *Workflow) InitCreatedAt() {
	action.CreatedAt = time.Now().UTC()
}

func (action *Workflow) InitUpdatedAt() {
	action.UpdatedAt = time.Now().UTC()
}

func (workflow *Workflow) UpdateWithRequest(userID uuid.UUID, req *request.UpdateWorkflowRequest) {
	workflow.ResourceID = idconvertor.ConvertStringToInt(req.ResourceID)
	workflow.Name = req.DisplayName
	workflow.Type = resourcelist.GetResourceNameMappedID(req.WorkflowType)
	workflow.TriggerMode = req.TriggerMode
	workflow.Transformer = req.ExportTransformerInString()
	workflow.Template = req.ExportTemplateInString()
	workflow.Config = req.ExportConfigInString()
	workflow.UpdatedBy = userID
	workflow.InitUpdatedAt()
}

func (action *Workflow) InitForFork(userID uuid.UUID) {
	action.CreatedBy = userID
	action.UpdatedBy = userID
	action.InitUID()
	action.InitCreatedAt()
	action.InitUpdatedAt()
}

func (action *Workflow) SetTemplate(template interface{}) {
	templateInJSONByte, _ := json.Marshal(template)
	action.Template = string(templateInJSONByte)
}

func (action *Workflow) ExportType() int {
	return action.Type
}

func (action *Workflow) ExportResourceID() int {
	return action.ResourceID
}

// func (action *Workflow) ExportConfig() *FlowActionConfig {
// 	ac := NewFlowActionConfig()
// 	json.Unmarshal([]byte(action.Config), ac)
// 	return ac
// }

func (action *Workflow) ExportDisplayName() string {
	return action.Name
}

func (action *Workflow) ExportIcon() string {
	content := action.ExportTemplateInMap()
	virtualResource, hitVirtualResource := content["virtualResource"]
	if !hitVirtualResource {
		return ""
	}
	virtualResourceAsserted, virtualResourceAssertPass := virtualResource.(map[string]interface{})
	if !virtualResourceAssertPass {
		return ""
	}
	icon, hitIcon := virtualResourceAsserted["icon"]
	if !hitIcon {
		return ""
	}
	iconAsserted, iconAssertPass := icon.(string)
	if !iconAssertPass {
		return ""
	}
	return iconAsserted
}

// func (action *Workflow) ExportTypeInString() string {
// 	return resourcelist.GetResourceIDMappedType(action.Type)
// }

// func (action *Workflow) SetContextByMap(context map[string]interface{}) {
// 	template := action.ExportTemplateInMap()
// 	template[ACTION_RUNTIME_INFO_FIELD_CONTEXT] = context
// 	templateJsonByte, _ := json.Marshal(template)
// 	action.Template = string(templateJsonByte)
// }

// func (action *Workflow) UpdateAppConfig(actionConfig *FlowActionConfig, userID string) {
// 	action.Config = actionConfig.ExportToJSONString()
// 	action.UpdatedBy = userID
// 	action.InitUpdatedAt()
// }

// func (action *Workflow) MergeRunFlowActionContextToRawTemplate(context map[string]interface{}) {
// 	template := action.ExportTemplateInMap()
// 	template[ACTION_RUNTIME_INFO_FIELD_CONTEXT] = context
// 	templateJsonByte, _ := json.Marshal(template)
// 	action.RawTemplate = string(templateJsonByte)
// }

// func (action *Workflow) UpdateWithRunFlowActionRequest(req *request.RunFlowActionRequest, userID string) {
// 	action.MergeRunFlowActionContextToRawTemplate(req.ExportContext())
// 	action.Template = req.ExportTemplateInString()
// 	action.UpdatedBy = userID
// 	action.InitUpdatedAt()

// 	// check if is onboarding action (which have no action storaged in database)
// 	if len(action.Template) == 0 {
// 		action.Template = action.RawTemplate
// 	}
// }

// func (action *Workflow) UpdateWorkflowByUpdateFlowActionRequest(teamID int, workflowID int, userID string, req *request.UpdateFlowActionRequest) {
// 	action.TeamID = teamID
// 	action.WorkflowID = workflowID
// 	// action.Version = APP_EDIT_VERSION // new action always created in builder edit mode, and it is edit version.
// 	action.ResourceID = idconvertor.ConvertStringToInt(req.ResourceID)
// 	action.Name = req.DisplayName
// 	// action.Type = resourcelist.GetResourceNameMappedID(req.FlowActionType)
// 	action.TriggerMode = req.TriggerMode
// 	action.Transformer = req.ExportTransformerInString()
// 	action.Template = req.ExportTemplateInString()
// 	action.Config = req.ExportConfigInString()
// 	action.UpdatedBy = userID
// 	action.InitUpdatedAt()
// }

// func (action *Workflow) IsVirtualFlowAction() bool {
// 	return resourcelist.IsVirtualResourceByIntType(action.Type)
// }

// func (action *Workflow) IsLocalVirtualFlowAction() bool {
// 	return resourcelist.IsLocalVirtualResourceByIntType(action.Type)
// }

// func (action *Workflow) IsRemoteVirtualFlowAction() bool {
// 	return resourcelist.IsRemoteVirtualResourceByIntType(action.Type)
// }

func (action *Workflow) ExportTransformerInMap() map[string]interface{} {
	var payload map[string]interface{}
	json.Unmarshal([]byte(action.Transformer), &payload)
	return payload
}

func (action *Workflow) ExportTemplateInMap() map[string]interface{} {
	payload := make(map[string]interface{}, 0)

	if action.Template != "" {
		json.Unmarshal([]byte(action.Template), &payload)
	}
	// add resourceID, runByAnonymous, teamID field for extend action runtime info
	payload["resourceID"] = action.ResourceID
	payload["runByAnonymous"] = true
	return payload
}

func (action *Workflow) ExportRawTemplateInMap() map[string]interface{} {
	var payload map[string]interface{}
	json.Unmarshal([]byte(action.RawTemplate), &payload)
	return payload
}

func (action *Workflow) ExportConfigInMap() map[string]interface{} {
	var payload map[string]interface{}
	json.Unmarshal([]byte(action.Config), &payload)
	return payload
}

// the action runtime does not pass the env info for virtual resource, so add them.
// func (action *Workflow) AppendRuntimeInfoForVirtualResource(authorization string, teamID int) {
// 	template := action.ExportTemplateInMap()
// 	template[ACTION_RUNTIME_INFO_FIELD_TEAM_ID] = teamID // the action.TeamID will invalied when onboarding
// 	template[ACTION_RUNTIME_INFO_FIELD_APP_ID] = action.WorkflowID
// 	template[ACTION_RUNTIME_INFO_FIELD_RESOURCE_ID] = action.ResourceID
// 	template[ACTION_RUNTIME_INFO_FIELD_ACTION_ID] = action.ID
// 	template[ACTION_RUNTIME_INFO_FIELD_AUTHORIZATION] = authorization
// 	template[ACTION_RUNTIME_INFO_FIELD_RUN_BY_ANONYMOUS] = (authorization == "")
// 	templateInByte, _ := json.Marshal(template)
// 	action.Template = string(templateInByte)
// }

// func (action *Workflow) SetResourceIDByAiAgent(aiAgent *illaresourcemanagersdk.AIAgentForExport) {
// 	action.ResourceID = aiAgent.ExportIDInInt()
// }

// func DoesFlowActionHasBeenCreated(actionID int) bool {
// 	return actionID > INVALIED_ACTION_ID
// }
