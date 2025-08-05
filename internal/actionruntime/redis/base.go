package redis

import (
	"crypto/tls"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

func (r *Connector) getConnectionWithOptions(resourceOptions map[string]interface{}) (*redis.Client, error) {
	if err := mapstructure.Decode(resourceOptions, &r.Resource); err != nil {
		return nil, err
	}

	options := redis.Options{
		Addr:     r.Resource.Host + ":" + r.Resource.Port,
		Username: r.Resource.DatabaseUsername,
		Password: r.Resource.DatabasePassword,
		DB:       r.Resource.DatabaseIndex,
	}
	if r.Resource.SSL {
		tlsConfig := tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: r.Resource.Host,
		}
		options.TLSConfig = &tlsConfig
	}
	rdb := redis.NewClient(&options)

	return rdb, nil
}
