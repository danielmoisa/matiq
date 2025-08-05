package couchdb

type resource struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Username string
	Password string
	SSL      bool
}

type action struct {
	Method   string `validate:"required,oneof=listRecords retrieveRecord createRecord updateRecord deleteRecord find getView"`
	Database string
	Opts     map[string]interface{}
}
