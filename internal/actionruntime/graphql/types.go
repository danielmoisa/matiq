package graphql

type Resource struct {
	BaseURL              string `validate:"required"`
	URLParams            []map[string]string
	Headers              []map[string]string
	Cookies              []map[string]string
	Authentication       string `validate:"required,oneof=none basic bearer apiKey"`
	AuthContent          map[string]string
	DisableIntrospection bool
}

type Action struct {
	Query     string
	Variables []map[string]interface{}
	Headers   []map[string]string
}
