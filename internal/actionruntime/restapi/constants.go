package restapi

const (
	METHOD_GET     = "GET"
	METHOD_POST    = "POST"
	METHOD_PUT     = "PUT"
	METHOD_DELETE  = "DELETE"
	METHOD_PATCH   = "PATCH"
	METHOD_HEAD    = "HEAD"
	METHOD_OPTIONS = "OPTIONS"

	BODY_NONE   = "none"
	BODY_RAW    = "raw"
	BODY_BINARY = "binary"
	BODY_FORM   = "form-data"
	BODY_XWFU   = "x-www-form-urlencoded"

	AUTH_NONE   = "none"
	AUTH_BASIC  = "basic"
	AUTH_BEARER = "bearer"
	AUTH_DIGEST = "digest"
	AUTH_OAUTH1 = "oauth1.0"
	AUTH_HAWK   = "hawk"
	AUTH_AWS    = "aws"

	VERIFY_MODE_SKIP = "skip"
	VERIFY_MODE_FULL = "verify-full"
	VERIFY_MODE_CA   = "verify-ca"
)
