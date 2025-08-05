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

func (controller *Controller) CreateFlowAction(c *gin.Context) {
	// fetch needed param
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// fetch payload
	createFlowActionRequest := request.NewCreateFlowActionRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&createFlowActionRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}
	fmt.Printf("createFlowActionRequest: %+v\n", createFlowActionRequest)

	// Parse user ID to UUID
	parsedUserID, errInParse := uuid.Parse(userID)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid user ID format: "+errInParse.Error())
		return
	}

	// init flowAction instance
	flowAction, errorInNewFlowAction := model.NewFlowActionByCreateRequest(parsedUserID, createFlowActionRequest)
	if errorInNewFlowAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_FLOW_ACTION, "error in create flowAction instance: "+errorInNewFlowAction.Error())
		return
	}
	fmt.Printf("flowAction: %+v\n", flowAction)

	// validate flowAction options
	// errInValidateActionOptions := controller.ValidateFlowActionTemplate(c, flowAction)
	// if errInValidateActionOptions != nil {
	// 	return
	// }

	// create flowAction
	_, errInCreateFlowAction := controller.Repository.FlowActionRepository.Create(flowAction)
	if errInCreateFlowAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_FLOW_ACTION, "create flowAction error: "+errInCreateFlowAction.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, response.NewCreateFlowActionResponse(flowAction))
}

func (controller *Controller) UpdateFlowAction(c *gin.Context) {
	// Get flowAction ID directly as string from URL parameter
	flowActionIDParam := c.Param("flowActionID")
	if flowActionIDParam == "" {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "flowAction ID is required")
		return
	}

	// Get user ID from Keycloak authentication
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// Parse flowAction ID string to UUID
	parsedFlowActionID, errInParse := uuid.Parse(flowActionIDParam)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid flow action ID format: "+errInParse.Error())
		return
	}

	// Parse user ID to UUID
	parsedUserID, errInParseUser := uuid.Parse(userID)
	if errInParseUser != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid user ID format: "+errInParseUser.Error())
		return
	}

	// Fetch the existing flowAction to verify ownership
	existingFlowAction, errInGetFlowAction := controller.Repository.FlowActionRepository.RetrieveFlowActionByID(parsedFlowActionID)
	if errInGetFlowAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_FLOW_ACTION, "get flowAction error: "+errInGetFlowAction.Error())
		return
	}

	// Check if the flowAction was created by the current user
	if existingFlowAction.CreatedBy != parsedUserID {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can only update flowActions that you created")
		return
	}

	// fetch payload
	updateFlowActionRequest := request.NewUpdateFlowActionRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&updateFlowActionRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}
	fmt.Printf("updateFlowActionRequest: %+v\n", updateFlowActionRequest)

	// Update the existing flowAction with new data
	existingFlowAction.UpdateWithRequest(parsedUserID, updateFlowActionRequest)
	fmt.Printf("updated flowAction: %+v\n", existingFlowAction)

	// validate flowAction options
	// errInValidateActionOptions := controller.ValidateFlowActionTemplate(c, existingFlowAction)
	// if errInValidateActionOptions != nil {
	// 	return
	// }

	// update flowAction in database
	errInUpdateFlowAction := controller.Repository.FlowActionRepository.Update(existingFlowAction)
	if errInUpdateFlowAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_FLOW_ACTION, "update flowAction error: "+errInUpdateFlowAction.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, response.NewUpdateFlowActionResponse(existingFlowAction))
}

func (controller *Controller) GetFlowAction(c *gin.Context) {
	// Get flowAction ID directly as string from URL parameter
	flowActionIDParam := c.Param("flowActionID")
	if flowActionIDParam == "" {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "flowAction ID is required")
		return
	}

	// Get user ID from Keycloak authentication
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// Parse flowAction ID string to UUID
	parsedFlowActionID, errInParse := uuid.Parse(flowActionIDParam)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid flowAction ID format: "+errInParse.Error())
		return
	}

	// Parse user ID to UUID
	parsedUserID, errInParse := uuid.Parse(userID)
	if errInParse != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_PARAM_FAILED, "invalid user ID format: "+errInParse.Error())
		return
	}

	// fetch data
	flowAction, errInGetAction := controller.Repository.FlowActionRepository.RetrieveFlowActionByID(parsedFlowActionID)
	if errInGetAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_FLOW_ACTION, "get flowAction error: "+errInGetAction.Error())
		return
	}

	// Check if the flowAction was created by the current user
	if flowAction.CreatedBy != parsedUserID {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can only access flowActions that you created")
		return
	}

	// new response
	getActionResponse := response.NewGetFlowActionResponse(flowAction)

	// feedback
	controller.FeedbackOK(c, getActionResponse)
}

func (controller *Controller) GetFlowActions(c *gin.Context) {
	// Get user ID from Keycloak authentication
	userID, errInGetUserID := controller.GetUserIDFromKeycloakAuth(c)
	if errInGetUserID != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "authentication required")
		return
	}

	// fetch data - this method already filters flowActions by user ID (created by user)
	flowActions, errInGetFlowActions := controller.Repository.FlowActionRepository.RetrieveAllFlowActionsByUserID(userID)
	if errInGetFlowActions != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_FLOW_ACTION, "get flowActions error: "+errInGetFlowActions.Error())
		return
	}

	// new response
	getFlowActionsResponse := response.NewGetFlowActionsResponse(flowActions)

	// feedback
	controller.FeedbackOK(c, getFlowActionsResponse)
}
