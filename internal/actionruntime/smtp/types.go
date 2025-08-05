package smtp

type Resource struct {
	Host     string `validate:"required"`
	Port     int    `validate:"gt=0"`
	Username string
	Password string
}

type Action struct {
	From        string   `validate:"required"`
	To          []string `validate:"required,gt=0,dive,required"`
	Bcc         []string
	Cc          []string
	SetReplyTo  bool
	ReplyTo     string `validate:"required_unless=SetReplyTo false"`
	Subject     string `validate:"required"`
	ContentType string `validate:"required,oneof=text/plain text/html"`
	Body        string `validate:"required"`
	Attachment  []Attachment
}

type Attachment struct {
	Data        string
	Name        string
	ContentType string
}
