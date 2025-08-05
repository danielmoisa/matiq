package model

import (
	"encoding/json"
	"time"

	"github.com/danielmoisa/matiq/internal/request"
	"github.com/danielmoisa/matiq/internal/utils/idconvertor"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/google/uuid"
)

type FlowAction struct {
	UID         uuid.UUID              `gorm:"column:uid;type:uuid;not null"`
	ResourceID  int                    `gorm:"column:resource_id;type:bigint;not null"`
	Name        string                 `gorm:"column:name;type:varchar;size:255;not null"`
	ActionType  int                    `gorm:"column:action_type;type:smallint;not null"`
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

func NewFlowAction() *FlowAction {
	return &FlowAction{}
}

func NewFlowActionByCreateRequest(userID uuid.UUID, req *request.CreateFlowActionRequest) (*FlowAction, error) {
	action := &FlowAction{
		ResourceID:  idconvertor.ConvertStringToInt(req.ResourceID),
		Name:        req.DisplayName,
		ActionType:  resourcelist.GetResourceNameMappedID(req.ActionType),
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

func NewFlowActionByUpdateRequest(userID uuid.UUID, req *request.UpdateFlowActionRequest) (*FlowAction, error) {
	action := &FlowAction{
		ResourceID:  idconvertor.ConvertStringToInt(req.ResourceID),
		Name:        req.DisplayName,
		ActionType:  resourcelist.GetResourceNameMappedID(req.ActionType),
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

func NewFlowActionByRunRequest(userID uuid.UUID, req *request.RunFlowActionRequest) *FlowAction {
	action := &FlowAction{
		ResourceID:  idconvertor.ConvertStringToInt(req.ResourceID),
		Name:        req.DisplayName,
		ActionType:  resourcelist.GetResourceNameMappedID(req.ActionType),
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

func (action *FlowAction) InitUID() {
	action.UID = uuid.New()
}

func (action *FlowAction) InitCreatedAt() {
	action.CreatedAt = time.Now().UTC()
}

func (action *FlowAction) InitUpdatedAt() {
	action.UpdatedAt = time.Now().UTC()
}

func (flowAction *FlowAction) UpdateWithRequest(userID uuid.UUID, req *request.UpdateFlowActionRequest) {
	flowAction.ResourceID = idconvertor.ConvertStringToInt(req.ResourceID)
	flowAction.Name = req.DisplayName
	flowAction.ActionType = resourcelist.GetResourceNameMappedID(req.ActionType)
	flowAction.TriggerMode = req.TriggerMode
	flowAction.Transformer = req.ExportTransformerInString()
	flowAction.Template = req.ExportTemplateInString()
	flowAction.Config = req.ExportConfigInString()
	flowAction.UpdatedBy = userID
	flowAction.InitUpdatedAt()
}

func (action *FlowAction) InitForFork(userID uuid.UUID) {
	action.CreatedBy = userID
	action.UpdatedBy = userID
	action.InitUID()
	action.InitCreatedAt()
	action.InitUpdatedAt()
}

func (action *FlowAction) SetTemplate(template interface{}) {
	templateInJSONByte, _ := json.Marshal(template)
	action.Template = string(templateInJSONByte)
}

func (action *FlowAction) ExportType() int {
	return action.ActionType
}

func (action *FlowAction) ExportResourceID() int {
	return action.ResourceID
}

func (action *FlowAction) ExportDisplayName() string {
	return action.Name
}

func (action *FlowAction) ExportIcon() string {
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

func (action *FlowAction) ExportTransformerInMap() map[string]interface{} {
	var payload map[string]interface{}
	json.Unmarshal([]byte(action.Transformer), &payload)
	return payload
}

func (action *FlowAction) ExportTemplateInMap() map[string]interface{} {
	payload := make(map[string]interface{}, 0)

	if action.Template != "" {
		json.Unmarshal([]byte(action.Template), &payload)
	}
	// add resourceID, runByAnonymous, teamID field for extend action runtime info
	payload["resourceID"] = action.ResourceID
	payload["runByAnonymous"] = true
	return payload
}

func (action *FlowAction) ExportRawTemplateInMap() map[string]interface{} {
	var payload map[string]interface{}
	json.Unmarshal([]byte(action.RawTemplate), &payload)
	return payload
}

func (action *FlowAction) ExportConfigInMap() map[string]interface{} {
	var payload map[string]interface{}
	json.Unmarshal([]byte(action.Config), &payload)
	return payload
}
