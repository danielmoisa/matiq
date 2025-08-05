package firebase

const (
	AUTH_SERVICE      = "auth"
	DATABASE_SERVICE  = "database"
	FIRESTORE_SERVICE = "firestore"
)

type Resource struct {
	DatabaseURL string                 `validate:"required,url"`
	ProjectID   string                 `validate:"required"`
	PrivateKey  map[string]interface{} `validate:"required"`
}

type Action struct {
	Service   string                 `validate:"required,oneof=firestore database auth"`
	Operation string                 `validate:"required"`
	Options   map[string]interface{} `validate:"required"`
}
