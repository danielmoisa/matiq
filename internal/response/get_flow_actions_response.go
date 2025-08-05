package response

import "github.com/danielmoisa/matiq/internal/model"

type GetFlowActionsResponse struct {
	Flows []*GetFlowActionResponse `json:"flows"`
}

func NewGetFlowActionsResponse(flows []*model.FlowAction) *GetFlowActionsResponse {
	var flowResponses []*GetFlowActionResponse
	for _, flow := range flows {
		flowResponses = append(flowResponses, NewGetFlowActionResponse(flow))
	}
	return &GetFlowActionsResponse{
		Flows: flowResponses,
	}
}

// ExportForFeedback implements the Response interface.
func (r *GetFlowActionsResponse) ExportForFeedback() interface{} {
	return r
}
