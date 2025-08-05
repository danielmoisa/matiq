package webhookresponse

import (
	"errors"
	"fmt"

	"github.com/danielmoisa/matiq/internal/actionruntime/common"
)

type WebhookResponseConnector struct {
	Action WebhookResponseTemplate
}

// server side transformer have no validate resource options method
func (r *WebhookResponseConnector) ValidateResourceOptions(resourceOptions map[string]interface{}) (common.ValidateResult, error) {
	return common.ValidateResult{Valid: true}, nil
}

func (r *WebhookResponseConnector) ValidateActionTemplate(actionOptions map[string]interface{}) (common.ValidateResult, error) {
	fmt.Printf("[DUMP] actionOptions: %+v \n", actionOptions)
	// @todo: check action needed field
	return common.ValidateResult{Valid: true}, nil
}

// server side transformer have no test connection method
func (r *WebhookResponseConnector) TestConnection(resourceOptions map[string]interface{}) (common.ConnectionResult, error) {
	return common.ConnectionResult{Success: false}, errors.New("unsupported type: server side transformer")
}

// server side transformer have no meta info
func (r *WebhookResponseConnector) GetMetaInfo(resourceOptions map[string]interface{}) (common.MetaInfoResult, error) {
	return common.MetaInfoResult{Success: false}, errors.New("unsupported type: server side transformer")
}

func (r *WebhookResponseConnector) Run(resourceOptions map[string]interface{}, actionOptions map[string]interface{}, rawActionOptions map[string]interface{}) (common.RuntimeResult, error) {
	res := common.RuntimeResult{
		Success: false,
		Rows:    []map[string]interface{}{},
		Extra:   map[string]interface{}{},
	}

	fmt.Printf("[DUMP] WebhookResponseConnector.Run() actionOptions: %+v\n", actionOptions)

	fmt.Printf("[DUMP] res: %+v\n", res)
	return res, nil
}
