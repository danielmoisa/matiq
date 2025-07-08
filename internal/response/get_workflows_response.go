package response

import "github.com/danielmoisa/workflow-builder/internal/model"

type GetWorkflowsResponse struct {
	Workflows []*GetWorkflowResponse `json:"workflows"`
}

func NewGetWorkflowsResponse(workflows []*model.Workflow) *GetWorkflowsResponse {
	var workflowResponses []*GetWorkflowResponse
	for _, workflow := range workflows {
		workflowResponses = append(workflowResponses, NewGetWorkflowResponse(workflow))
	}
	return &GetWorkflowsResponse{
		Workflows: workflowResponses,
	}
}

// ExportForFeedback implements the Response interface.
func (r *GetWorkflowsResponse) ExportForFeedback() interface{} {
	return r
}
