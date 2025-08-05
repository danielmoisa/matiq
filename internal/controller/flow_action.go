package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) CreateFlowAction(c *gin.Context) {
	// fetch needed param
	//teamID, errInGetTeamID := controller.GetMagicIntParamFromRequest(c, PARAM_TEAM_ID)
	//workflowID, errInGetWORKFLOWID := controller.GetMagicIntParamFromRequest(c, PARAM_WORKFLOW_ID)
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	userAuthToken, errInGetAuthToken := controller.GetUserAuthTokenFromHeader(c)
	if errInGetTeamID != nil || errInGetWORKFLOWID != nil || errInGetUserID != nil || errInGetAuthToken != nil {
		return
	}

	// validate
	canManage, errInCheckAttr := controller.AttributeGroup.CanManage(
		teamID,
		userAuthToken,
		accesscontrol.UNIT_TYPE_FLOW_ACTION,
		accesscontrol.DEFAULT_UNIT_ID,
		accesscontrol.ACTION_MANAGE_CREATE_FLOW_ACTION,
	)
	if errInCheckAttr != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
		return
	}
	if !canManage {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// fetch payload
	createFlowActionRequest := request.NewCreateFlowActionRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&createFlowActionRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}
	fmt.Printf("createFlowActionRequest: %+v\n", createFlowActionRequest)

	// append remote virtual resource (like aiagent, but the transformet is local virtual resource)
	if createFlowActionRequest.IsRemoteVirtualAction() {
		// the AI_Agent need fetch resource info from resource manager, but illa drive does not need that
		if createFlowActionRequest.NeedFetchResourceInfoFromSourceManager() {
			api, errInNewAPI := illaresourcemanagersdk.NewIllaResourceManagerRestAPI()
			if errInNewAPI != nil {
				controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_FLOW_ACTION, "error in fetch flowAction mapped virtual resource: "+errInNewAPI.Error())
				return
			}
			virtualResource, errInGetVirtualResource := api.GetResource(createFlowActionRequest.ExportFlowActionTypeInInt(), createFlowActionRequest.ExportResourceIDInInt())
			if errInGetVirtualResource != nil {
				controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_FLOW_ACTION, "error in fetch flowAction mapped virtual resource: "+errInGetVirtualResource.Error())
				return
			}
			createFlowActionRequest.AppendVirtualResourceToTemplate(virtualResource)
		}
	}

}
