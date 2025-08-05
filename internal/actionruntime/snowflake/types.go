package snowflake

import (
	"errors"

	"github.com/danielmoisa/matiq/internal/actionruntime/common"
)

const (
	FIELD_CONTEXT = "context"
	FIELD_QUERY   = "query"
)

type Resource struct {
	AccountName    string `validate:"required"`
	Warehouse      string `validate:"required"`
	Database       string `validate:"required"`
	Schema         string
	Role           string
	Authentication string            `validate:"oneof=basic key"`
	AuthContent    map[string]string `validate:"required"`
}

type Action struct {
	Mode     string `validate:"oneof=gui sql sql-safe"`
	Query    string
	RawQuery string
	Context  map[string]interface{}
}

func (q *Action) IsSafeMode() bool {
	return q.Mode == common.MODE_SQL_SAFE
}

func (q *Action) SetRawQueryAndContext(rawTemplate map[string]interface{}) error {
	queryRaw, hit := rawTemplate[FIELD_QUERY]
	if !hit {
		return errors.New("missing query field for SetRawQueryAndContext() in query")
	}
	queryAsserted, assertPass := queryRaw.(string)
	if !assertPass {
		return errors.New("query field assert failed in SetRawQueryAndContext() method")

	}
	q.RawQuery = queryAsserted
	contextRaw, hit := rawTemplate[FIELD_CONTEXT]
	if !hit {
		return errors.New("missing context field SetRawQueryAndContext() in query")
	}
	contextAsserted, assertPass := contextRaw.(map[string]interface{})
	if !assertPass {
		return errors.New("context field assert failed in SetRawQueryAndContext() method")

	}
	q.Context = contextAsserted
	return nil
}
