package trigger

import (
	"errors"
	"fmt"

	"github.com/danielmoisa/matiq/internal/actionruntime/common"
)

type TriggerConnector struct {
	Action TriggerTemplate
}

// AI Agent have no validate resource options method
func (r *TriggerConnector) ValidateResourceOptions(resourceOptions map[string]interface{}) (common.ValidateResult, error) {
	return common.ValidateResult{Valid: true}, nil
}

func (r *TriggerConnector) ValidateActionTemplate(actionOptions map[string]interface{}) (common.ValidateResult, error) {
	fmt.Printf("[DUMP] actionOptions: %+v \n", actionOptions)
	// @todo: check action needed field
	return common.ValidateResult{Valid: true}, nil
}

// AI Agent have no test connection method
func (r *TriggerConnector) TestConnection(resourceOptions map[string]interface{}) (common.ConnectionResult, error) {
	return common.ConnectionResult{Success: false}, errors.New("unsupported type: AI Agent")
}

// AI Agent have no meta info
func (r *TriggerConnector) GetMetaInfo(resourceOptions map[string]interface{}) (common.MetaInfoResult, error) {
	return common.MetaInfoResult{Success: false}, errors.New("unsupported type: AI Agent")
}

func (r *TriggerConnector) Run(resourceOptions map[string]interface{}, actionOptions map[string]interface{}, rawActionOptions map[string]interface{}) (common.RuntimeResult, error) {
	res := common.RuntimeResult{
		Success: false,
		Rows:    []map[string]interface{}{},
		Extra:   map[string]interface{}{},
	}

	fmt.Printf("[DUMP] illadrive.Run() actionOptions: %+v\n", actionOptions)

	fmt.Printf("[DUMP] res: %+v\n", res)
	return res, nil
}
