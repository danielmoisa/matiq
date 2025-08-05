package firebase

import (
	"context"
	"encoding/json"

	firebase "firebase.google.com/go/v4"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
)

func (f *Connector) getConnectionWithOptions(resourceOptions map[string]interface{}) (*firebase.App, error) {
	if err := mapstructure.Decode(resourceOptions, &f.ResourceOpts); err != nil {
		return nil, err
	}

	// build firebase service account
	privateKey, err := json.Marshal(f.ResourceOpts.PrivateKey)
	if err != nil {
		return nil, err
	}
	sa := option.WithCredentialsJSON(privateKey)

	// build firebase config for realtime database
	config := &firebase.Config{DatabaseURL: f.ResourceOpts.DatabaseURL}

	// new firebase app
	app, err := firebase.NewApp(context.Background(), config, sa)
	if err != nil {
		return nil, err
	}

	return app, nil
}
