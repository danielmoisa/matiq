package elasticsearch

const (
	SEARCH_OPERATION = "search"
	INSERT_OPERATION = "insert"
	GET_OPERATION    = "get"
	UPDATE_OPERATION = "update"
	DELETE_OPERATION = "delete"
)

type Resource struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type Action struct {
	Operation string `validate:"required,oneof=search insert get update delete"`
	Index     string
	ID        string
	Body      string
	Query     string
}
