package common

const (
	MODE_GUI      = "gui"
	MODE_SQL      = "sql"
	MODE_SQL_SAFE = "sql-safe"
)

type ValidateResult struct {
	Valid bool
	Extra map[string]interface{}
}

type ConnectionResult struct {
	Success bool
}

type RuntimeResult struct {
	Success bool
	Rows    []map[string]interface{}
	Extra   map[string]interface{}
}

func (i *RuntimeResult) SetSuccess() {
	i.Success = true
}

type MetaInfoResult struct {
	Success bool
	Schema  map[string]interface{}
}

func (metaInfoResult *MetaInfoResult) ExportSchema() map[string]interface{} {
	return metaInfoResult.Schema
}
