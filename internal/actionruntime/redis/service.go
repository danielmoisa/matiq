package redis

import (
	"context"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"

	"github.com/danielmoisa/matiq/internal/actionruntime/common"
)

type Connector struct {
	Resource Options
	Action   Command
}

func (r *Connector) ValidateResourceOptions(resourceOptions map[string]interface{}) (common.ValidateResult, error) {
	// format resource options
	if err := mapstructure.Decode(resourceOptions, &r.Resource); err != nil {
		return common.ValidateResult{Valid: false}, err
	}

	// validate redis options
	validate := validator.New()
	if err := validate.Struct(r.Resource); err != nil {
		return common.ValidateResult{Valid: false}, err
	}
	return common.ValidateResult{Valid: true}, nil
}

func (r *Connector) ValidateActionTemplate(actionOptions map[string]interface{}) (common.ValidateResult, error) {
	// format redis command options
	if err := mapstructure.Decode(actionOptions, &r.Action); err != nil {
		return common.ValidateResult{Valid: false}, err
	}

	// validate redis command options
	validate := validator.New()
	if err := validate.Struct(r.Action); err != nil {
		return common.ValidateResult{Valid: false}, err
	}
	return common.ValidateResult{Valid: true}, nil
}

func (r *Connector) TestConnection(resourceOptions map[string]interface{}) (common.ConnectionResult, error) {
	// get redis client
	rdb, err := r.getConnectionWithOptions(resourceOptions)
	if err != nil {
		return common.ConnectionResult{Success: false}, err
	}
	defer rdb.Close()

	// test redis connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return common.ConnectionResult{Success: false}, err
	}

	return common.ConnectionResult{Success: true}, nil
}

func (r *Connector) GetMetaInfo(resourceOptions map[string]interface{}) (common.MetaInfoResult, error) {

	return common.MetaInfoResult{
		Success: true,
		Schema:  nil,
	}, nil
}

func (r *Connector) Run(resourceOptions map[string]interface{}, actionOptions map[string]interface{}, rawActionOptions map[string]interface{}) (common.RuntimeResult, error) {
	// start a default context
	ctx, cancel := context.WithTimeout(context.TODO(), common.DEFAULT_QUERY_AND_EXEC_TIMEOUT)
	defer cancel()

	// get redis connection
	rdb, err := r.getConnectionWithOptions(resourceOptions)
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	defer rdb.Close()

	// test redis connection
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// format redis command
	if err := mapstructure.Decode(actionOptions, &r.Action); err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	redisCMD := strings.TrimSpace(r.Action.Query)
	redisCMDSlice := strings.Fields(redisCMD)
	inputRedisCMDSlice := make([]interface{}, len(redisCMDSlice))
	for i, v := range redisCMDSlice {
		inputRedisCMDSlice[i] = v
	}

	// run redis command
	val, err := rdb.Do(ctx, inputRedisCMDSlice...).Result()
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	cmdResult := common.RuntimeResult{
		Success: true,
		Rows:    []map[string]interface{}{},
		Extra:   map[string]interface{}{},
	}
	cmdResult.Rows = append(cmdResult.Rows, map[string]interface{}{"result": val})

	return cmdResult, nil
}
