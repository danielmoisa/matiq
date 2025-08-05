package couchdb

import (
	"fmt"
	"net/url"

	_ "github.com/go-kivik/couchdb/v4"
	"github.com/go-kivik/kivik/v4"
	"github.com/mitchellh/mapstructure"
)

const (
	LIST_METHOD     = "listRecords"
	RETRIEVE_METHOD = "retrieveRecord"
	CREATE_METHOD   = "createRecord"
	UPDATE_METHOD   = "updateRecord"
	DELETE_METHOD   = "deleteRecord"
	FIND_METHOD     = "find"
	GET_METHOD      = "getView"
)

func (c *Connector) getClient(resourceOptions map[string]interface{}) (*kivik.Client, error) {
	// format resource options
	if err := mapstructure.Decode(resourceOptions, &c.resourceOptions); err != nil {
		return nil, err
	}

	protocolStr := "http"
	if c.resourceOptions.SSL {
		protocolStr = "https"
	}
	escapedPassword := url.QueryEscape(c.resourceOptions.Password)
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/", protocolStr, c.resourceOptions.Username, escapedPassword,
		c.resourceOptions.Host, c.resourceOptions.Port)
	client, err := kivik.New("couch", dsn)
	if err != nil {
		return nil, err
	}

	return client, nil
}
