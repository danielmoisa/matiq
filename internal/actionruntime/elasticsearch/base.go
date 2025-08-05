package elasticsearch

import (
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
)

func (e *Connector) getConnectionWithOptions(resourceOptions map[string]interface{}) (*es.Client, error) {
	if err := mapstructure.Decode(resourceOptions, &e.ResourceOpts); err != nil {
		return nil, err
	}

	esCfg := es.Config{
		Addresses: []string{
			e.ResourceOpts.Host + ":" + e.ResourceOpts.Port,
		},
		Username: e.ResourceOpts.Username,
		Password: e.ResourceOpts.Password,
	}
	esClient, err := es.NewClient(esCfg)
	if err != nil {
		return nil, err
	}
	return esClient, err
}
