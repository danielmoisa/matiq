package redis

type Options struct {
	Host             string `validate:"required"`
	Port             string `validate:"required"`
	DatabaseIndex    int    `validate:"gte=0"`
	DatabaseUsername string
	DatabasePassword string
	SSL              bool
}

type Command struct {
	Mode  string `validate:"required,oneof=select raw"`
	Query string
}
