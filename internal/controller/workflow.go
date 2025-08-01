package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danielmoisa/matiq/internal/model"
	"github.com/danielmoisa/matiq/internal/request"
	"github.com/danielmoisa/matiq/internal/response"
)

func (controller *Controller) CreateWorkflow(c *gin.Context) {
	// fetch needed param
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// fetch payload
	createWorkflowRequest := request.NewCreateWorkflowRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&createWorkflowRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}
	fmt.Printf("createWorkflowRequest: %+v\n", createWorkflowRequest)

	// append remote virtual resource (like aiagent, but the transformet is local virtual resource)
	// if createFlowActionRequest.IsRemoteVirtualAction() {
	// 	// the AI_Agent need fetch resource info from resource manager, but illa drive does not need that
	// 	if createFlowActionRequest.NeedFetchResourceInfoFromSourceManager() {
	// 		api, errInNewAPI := illaresourcemanagersdk.NewIllaResourceManagerRestAPI()
	// 		if errInNewAPI != nil {
	// 			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_FLOW_ACTION, "error in fetch flowAction mapped virtual resource: "+errInNewAPI.Error())
	// 			return
	// 		}
	// 		virtualResource, errInGetVirtualResource := api.GetResource(createFlowActionRequest.ExportFlowActionTypeInInt(), createFlowActionRequest.ExportResourceIDInInt())
	// 		if errInGetVirtualResource != nil {
	// 			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_FLOW_ACTION, "error in fetch flowAction mapped virtual resource: "+errInGetVirtualResource.Error())
	// 			return
	// 		}
	// 		createFlowActionRequest.AppendVirtualResourceToTemplate(virtualResource)
	// 	}
	// }

	// Parse user ID to UUID
	parsedUserID, errInParse := uuid.Parse(userID)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid user ID format: "+errInParse.Error())
		return
	}

	// init workflow instace
	workflow, errorInNewWorkflow := model.NewWorkflowByCreateRequest(parsedUserID, createWorkflowRequest)
	if errorInNewWorkflow != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_WORKFLOW, "error in create workflow instance: "+errorInNewWorkflow.Error())
		return
	}
	fmt.Printf("workflow: %+v\n", workflow)

	// validate flowAction options
	// errInValidateActionOptions := controller.ValidateFlowActionTemplate(c, flowAction)
	// if errInValidateActionOptions != nil {
	// 	return
	// }

	// create workflow
	_, errInCreateWorkflow := controller.Repository.WorkflowRepository.Create(workflow)
	if errInCreateWorkflow != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_WORKFLOW, "create workflow error: "+errInCreateWorkflow.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, response.NewCreateWorkflowResponse(workflow))
}

func (controller *Controller) UpdateWorkflow(c *gin.Context) {
	// Get workflow ID directly as string from URL parameter
	workflowIDParam := c.Param("workflowID")
	if workflowIDParam == "" {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "workflow ID is required")
		return
	}

	// Get user ID from Keycloak authentication
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// Parse workflow ID string to UUID
	parsedWorkflowID, errInParse := uuid.Parse(workflowIDParam)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid workflow ID format: "+errInParse.Error())
		return
	}

	// Parse user ID to UUID
	parsedUserID, errInParseUser := uuid.Parse(userID)
	if errInParseUser != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid user ID format: "+errInParseUser.Error())
		return
	}

	// Fetch the existing workflow to verify ownership
	existingWorkflow, errInGetWorkflow := controller.Repository.WorkflowRepository.RetrieveWorkflowByID(parsedWorkflowID)
	if errInGetWorkflow != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_WORKFLOW, "get workflow error: "+errInGetWorkflow.Error())
		return
	}

	// Check if the workflow was created by the current user
	if existingWorkflow.CreatedBy != parsedUserID {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can only update workflows that you created")
		return
	}

	// fetch payload
	updateWorkflowRequest := request.NewUpdateWorkflowRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&updateWorkflowRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}
	fmt.Printf("updateWorkflowRequest: %+v\n", updateWorkflowRequest)

	// Update the existing workflow with new data
	existingWorkflow.UpdateWithRequest(parsedUserID, updateWorkflowRequest)
	fmt.Printf("updated workflow: %+v\n", existingWorkflow)

	// validate workflow options
	// errInValidateActionOptions := controller.ValidateFlowActionTemplate(c, existingWorkflow)
	// if errInValidateActionOptions != nil {
	// 	return
	// }

	// update workflow in database
	errInUpdateWorkflow := controller.Repository.WorkflowRepository.Update(existingWorkflow)
	if errInUpdateWorkflow != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_WORKFLOW, "update workflow error: "+errInUpdateWorkflow.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, response.NewUpdateWorkflowResponse(existingWorkflow))
}

// func (controller *Controller) UpdateFlowAction(c *gin.Context) {
// 	// fetch needed param
// 	teamID, errInGetTeamID := controller.GetMagicIntParamFromRequest(c, PARAM_TEAM_ID)
// 	workflowID, errInGetWORKFLOWID := controller.GetMagicIntParamFromRequest(c, PARAM_WORKFLOW_ID)
// 	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
// 	flowActionID, errInGetActionID := controller.GetMagicIntParamFromRequest(c, PARAM_FLOW_ACTION_ID)
// 	userAuthToken, errInGetAuthToken := controller.GetUserAuthTokenFromHeader(c)
// 	if errInGetTeamID != nil || errInGetWORKFLOWID != nil || errInGetUserID != nil || errInGetActionID != nil || errInGetAuthToken != nil {
// 		return
// 	}

// 	// validate
// 	canManage, errInCheckAttr := controller.AttributeGroup.CanManage(
// 		teamID,
// 		userAuthToken,
// 		accesscontrol.UNIT_TYPE_FLOW_ACTION,
// 		flowActionID,
// 		accesscontrol.ACTION_MANAGE_EDIT_FLOW_ACTION,
// 	)
// 	if errInCheckAttr != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
// 		return
// 	}
// 	if !canManage {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
// 		return
// 	}

// 	// fetch payload
// 	updateFlowActionRequest := request.NewUpdateFlowActionRequest()
// 	if err := json.NewDecoder(c.Request.Body).Decode(&updateFlowActionRequest); err != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
// 		return
// 	}

// 	// update
// 	newInDatabaseFlowAction, errInUpdateAction := controller.updateFlowAction(c, teamID, workflowID, userID, flowActionID, updateFlowActionRequest)
// 	if errInUpdateAction != nil {
// 		return
// 	}

// 	// feedback
// 	controller.FeedbackOK(c, response.NewUpdateFlowActionResponse(newInDatabaseFlowAction))
// }

// func (controller *Controller) UpdateFlowActionByBatch(c *gin.Context) {
// 	// fetch needed param
// 	teamID, errInGetTeamID := controller.GetMagicIntParamFromRequest(c, PARAM_TEAM_ID)
// 	workflowID, errInGetAPPID := controller.GetMagicIntParamFromRequest(c, PARAM_WORKFLOW_ID)
// 	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
// 	userAuthToken, errInGetAuthToken := controller.GetUserAuthTokenFromHeader(c)
// 	if errInGetTeamID != nil || errInGetAPPID != nil || errInGetUserID != nil || errInGetAuthToken != nil {
// 		return
// 	}

// 	// validate
// 	canManage, errInCheckAttr := controller.AttributeGroup.CanManage(
// 		teamID,
// 		userAuthToken,
// 		accesscontrol.UNIT_TYPE_ACTION,
// 		0,
// 		accesscontrol.ACTION_MANAGE_EDIT_ACTION,
// 	)
// 	if errInCheckAttr != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
// 		return
// 	}
// 	if !canManage {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
// 		return
// 	}

// 	// fetch payload
// 	updateFlowActionByBatchRequest := request.NewUpdateFlowActionByBatchRequest()
// 	if err := json.NewDecoder(c.Request.Body).Decode(&updateFlowActionByBatchRequest); err != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
// 		return
// 	}

// 	inDatabaseFlowActions := make([]*model.FlowAction, 0)
// 	for _, updateFlowActionRequest := range updateFlowActionByBatchRequest.ExportFlowActions() {
// 		newInDatabaseFlowAction, errInUpdateAction := controller.updateFlowAction(c, teamID, workflowID, userID, updateFlowActionRequest.ExportFlowActionIDInInt(), updateFlowActionRequest)
// 		if errInUpdateAction != nil {
// 			return
// 		}
// 		inDatabaseFlowActions = append(inDatabaseFlowActions, newInDatabaseFlowAction)
// 	}

// 	// feedback
// 	controller.FeedbackOK(c, response.NewUpdateFlowActionByBatchResponse(inDatabaseFlowActions))
// }

// func (controller *Controller) DeleteFlowAction(c *gin.Context) {
// 	// fetch needed param
// 	teamID, errInGetTeamID := controller.GetMagicIntParamFromRequest(c, PARAM_TEAM_ID)
// 	flowActionID, errInGetActionID := controller.GetMagicIntParamFromRequest(c, PARAM_FLOW_ACTION_ID)
// 	userAuthToken, errInGetAuthToken := controller.GetUserAuthTokenFromHeader(c)
// 	if errInGetTeamID != nil || errInGetActionID != nil || errInGetAuthToken != nil {
// 		return
// 	}

// 	// validate
// 	canManage, errInCheckAttr := controller.AttributeGroup.CanDelete(
// 		teamID,
// 		userAuthToken,
// 		accesscontrol.UNIT_TYPE_FLOW_ACTION,
// 		flowActionID,
// 		accesscontrol.ACTION_DELETE,
// 	)
// 	if errInCheckAttr != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
// 		return
// 	}
// 	if !canManage {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
// 		return
// 	}

// 	// delete
// 	errInDelete := controller.Storage.FlowActionStorage.DeleteFlowActionByTeamIDAndFlowActionID(teamID, flowActionID)
// 	if errInDelete != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_DELETE_FLOW_ACTION, "delete flowAction error: "+errInDelete.Error())
// 		return
// 	}

// 	// feedback
// 	controller.FeedbackOK(c, response.NewDeleteFlowActionResponse(flowActionID))
// 	return
// }

func (controller *Controller) GetWorkflow(c *gin.Context) {
	// Get workflow ID directly as string from URL parameter
	workflowIDParam := c.Param("workflowID")
	if workflowIDParam == "" {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "workflow ID is required")
		return
	}

	// Get user ID from Keycloak authentication
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// Parse workflow ID string to UUID
	parsedWorkflowID, errInParse := uuid.Parse(workflowIDParam)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid workflow ID format: "+errInParse.Error())
		return
	}

	// Parse user ID to UUID
	parsedUserID, errInParse := uuid.Parse(userID)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid user ID format: "+errInParse.Error())
		return
	}

	// fetch data
	workflow, errInGetAction := controller.Repository.WorkflowRepository.RetrieveWorkflowByID(parsedWorkflowID)
	if errInGetAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_WORKFLOW, "get workflow error: "+errInGetAction.Error())
		return
	}

	// Check if the workflow was created by the current user
	if workflow.CreatedBy != parsedUserID {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can only access workflows that you created")
		return
	}

	// new response
	getActionResponse := response.NewGetWorkflowResponse(workflow)

	// append remote virtual resource
	// if flowAction.IsRemoteVirtualFlowAction() {
	// 	api, errInNewAPI := illaresourcemanagersdk.NewIllaResourceManagerRestAPI()
	// 	if errInNewAPI != nil {
	// 		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_FLOW_ACTION, "error in fetch flowAction mapped virtual resource: "+errInNewAPI.Error())
	// 		return
	// 	}
	// 	virtualResource, errInGetVirtualResource := api.GetResource(flowAction.ExportType(), flowAction.ExportResourceID())
	// 	if errInGetVirtualResource != nil {
	// 		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_FLOW_ACTION, "error in fetch flowAction mapped virtual resource: "+errInGetVirtualResource.Error())
	// 		return
	// 	}
	// 	getActionResponse.AppendVirtualResourceToTemplate(virtualResource)
	// }

	// feedback
	controller.FeedbackOK(c, getActionResponse)
}

func (controller *Controller) GetWorkflows(c *gin.Context) {
	// Get user ID from Keycloak authentication
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// fetch data - this method already filters workflows by user ID (created by user)
	workflows, errInGetWorkflows := controller.Repository.WorkflowRepository.RetrieveAllWorkflowByUserID(userID)
	if errInGetWorkflows != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_WORKFLOW, "get workflows error: "+errInGetWorkflows.Error())
		return
	}

	// new response
	getWorkflowsResponse := response.NewGetWorkflowsResponse(workflows)

	// feedback
	controller.FeedbackOK(c, getWorkflowsResponse)
}

// func (controller *Controller) RunFlowAction(c *gin.Context) {
// 	// fetch needed param
// 	teamID, errInGetTeamID := controller.GetMagicIntParamFromRequest(c, PARAM_TEAM_ID)
// 	workflowID, errInGetAppID := controller.GetMagicIntParamFromRequest(c, PARAM_WORKFLOW_ID)
// 	flowActionID, errInGetActionID := controller.GetMagicIntParamFromRequest(c, PARAM_FLOW_ACTION_ID)
// 	userAuthToken, errInGetAuthToken := controller.GetUserAuthTokenFromHeader(c)
// 	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
// 	if errInGetTeamID != nil || errInGetAppID != nil || errInGetActionID != nil || errInGetAuthToken != nil || errInGetUserID != nil {
// 		return
// 	}

// 	// validate
// 	canManage, errInCheckAttr := controller.AttributeGroup.CanManage(
// 		teamID,
// 		userAuthToken,
// 		accesscontrol.UNIT_TYPE_FLOW_ACTION,
// 		flowActionID,
// 		accesscontrol.ACTION_MANAGE_RUN_FLOW_ACTION,
// 	)
// 	if errInCheckAttr != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
// 		return
// 	}
// 	if !canManage {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
// 		return
// 	}

// 	// set resource timing header
// 	// @see:
// 	// [Timing-Allow-Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Timing-Allow-Origin)
// 	// [Resource_timing](https://developer.mozilla.org/en-US/docs/Web/API/Performance_API/Resource_timing)
// 	c.Header("Timing-Allow-Origin", "*")

// 	// execute
// 	runFlowActionRequest := request.NewRunFlowActionRequest()
// 	if err := json.NewDecoder(c.Request.Body).Decode(&runFlowActionRequest); err != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error"+err.Error())
// 		return
// 	}

// 	// get flowAction
// 	flowAction := model.NewFlowAction()
// 	fmt.Printf("[RetrieveActionsByTeamIDActionID] teamID: %d, flowActionID: %d\n", teamID, flowActionID)

// 	// flowActionID has not been created (like flowActionID is 0 'ILAfx4p1C7d0'), but we still can run it (onboarding case)
// 	if !model.DoesActionHasBeenCreated(flowActionID) {
// 		// ok, flowAction was not created, fetch app and build a temporary flowAction.
// 		flowAction = model.NewFlowAcitonByRunFlowActionRequest(teamID, workflowID, userID, runFlowActionRequest)
// 	} else {
// 		// ok, we retrieve flowAction from database
// 		var errInRetrieveAction error
// 		flowAction, errInRetrieveAction = controller.Storage.FlowActionStorage.RetrieveFlowActionByTeamIDFlowActionID(teamID, flowActionID)
// 		if errInRetrieveAction != nil {
// 			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_FLOW_ACTION, "get flowAction failed: "+errInRetrieveAction.Error())
// 			return
// 		}
// 	}

// 	// update flowAction data with run flowAction reqeust
// 	flowAction.UpdateWithRunFlowActionRequest(runFlowActionRequest, userID)
// 	fmt.Printf("[DUMP] flowAction: %+v\n", flowAction)

// 	// process input context with action template
// 	// @todo: this method should rewrite to common method for all flow actions.
// 	avaliableDoProcessList := map[int]bool{
// 		resourcelist.TYPE_MONGODB_ID:  true,
// 		resourcelist.TYPE_APPWRITE_ID: true,
// 		resourcelist.TYPE_AIRTABLE_ID: true,
// 	}
// 	if avaliableDoProcessList[flowAction.ExportType()] {
// 		fmt.Printf("[DUMP] flowAction.ExportTemplateInMap() original: %+v\n", flowAction.ExportTemplateInMap())
// 		processedTemplate, errInProcessTemplate := common.ProcessTemplateByContext(flowAction.ExportTemplateInMap(), runFlowActionRequest.ExportContext())
// 		if errInProcessTemplate != nil {
// 			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_PROCESS_FLOW_ACTION, "process flow action failed: "+errInProcessTemplate.Error())
// 			return
// 		}
// 		flowAction.SetTemplate(processedTemplate)
// 		fmt.Printf("[DUMP] flowAction.ExportTemplateInMap() converted: %+v\n", flowAction.ExportTemplateInMap())
// 		processedTemplateInJSONbyte, _ := json.Marshal(processedTemplate)
// 		fmt.Printf("[DUMP] flowAction.ExportTemplateInMap() converted in json: %+v\n", string(processedTemplateInJSONbyte))
// 	}

// 	// assembly flowAction
// 	flowActionFactory := model.NewFlowActionFactoryByFlowAction(flowAction)
// 	flowActionAssemblyLine, errInBuild := flowActionFactory.Build()
// 	if errInBuild != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate flowAction type error: "+errInBuild.Error())
// 		return
// 	}

// 	// get resource
// 	resource := model.NewResource()
// 	if !flowAction.IsVirtualFlowAction() {
// 		// process normal resource flowAction
// 		var errInRetrieveResource error
// 		resource, errInRetrieveResource = controller.Storage.ResourceStorage.RetrieveByTeamIDAndResourceID(teamID, flowAction.ExportResourceID())
// 		if errInRetrieveResource != nil {
// 			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_RESOURCE, "get resource failed: "+errInRetrieveResource.Error())
// 			return
// 		}
// 		// resource option validate only happend in create or update phrase
// 		// note that validate will set resprce options to flowActionAssemblyLine
// 		_, errInValidateResourceOptions := flowActionAssemblyLine.ValidateResourceOptions(resource.ExportOptionsInMap())
// 		if errInValidateResourceOptions != nil {
// 			controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_RESOURCE_FAILED, "validate resource failed: "+errInValidateResourceOptions.Error())
// 			return
// 		}
// 	} else {
// 		// process virtual resource flowAction
// 		flowAction.AppendRuntimeInfoForVirtualResource(userAuthToken, teamID)
// 	}

// 	// check flowAction template
// 	fmt.Printf("[DUMP] flowAction.ExportTemplateInMap(): %+v\n", flowAction.ExportTemplateInMap())
// 	fmt.Printf("[DUMP] flowAction.ExportRawTemplateInMap(): %+v\n", flowAction.ExportRawTemplateInMap())
// 	_, errInValidate := flowActionAssemblyLine.ValidateActionTemplate(flowAction.ExportTemplateInMap())
// 	if errInValidate != nil {
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate flowAction template error: "+errInValidate.Error())
// 		return
// 	}

// 	// run
// 	log.Printf("[DUMP]flowAction: %+v\n", flowAction)
// 	log.Printf("[DUMP] resource.ExportOptionsInMap(): %+v, flowAction.ExportTemplateInMap(): %+v\n", resource.ExportOptionsInMap(), flowAction.ExportTemplateInMap())
// 	flowActionRunResult, errInRunAction := flowActionAssemblyLine.Run(resource.ExportOptionsInMap(), flowAction.ExportTemplateInMap(), flowAction.ExportRawTemplateInMap())
// 	if errInRunAction != nil {
// 		if strings.HasPrefix(errInRunAction.Error(), "Error 1064:") {
// 			lineNumber, _ := strconv.Atoi(errInRunAction.Error()[len(errInRunAction.Error())-1:])
// 			message := ""
// 			regexp, _ := regexp.Compile(`to use`)
// 			match := regexp.FindStringIndex(errInRunAction.Error())
// 			if len(match) == 2 {
// 				message = errInRunAction.Error()[match[1]:]
// 			}
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"errorCode":    400,
// 				"errorFlag":    ERROR_FLAG_EXECUTE_FLOW_ACTION_FAILED,
// 				"errorMessage": errors.New("SQL syntax error").Error(),
// 				"errorData": map[string]interface{}{
// 					"lineNumber": lineNumber,
// 					"message":    "SQL syntax error" + message,
// 				},
// 			})
// 			return
// 		}
// 		controller.FeedbackBadRequest(c, ERROR_FLAG_EXECUTE_FLOW_ACTION_FAILED, "run flowAction error: "+errInRunAction.Error())
// 		return
// 	}

// 	// feedback
// 	c.JSON(http.StatusOK, flowActionRunResult)
// }
