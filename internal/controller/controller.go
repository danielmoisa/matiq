package controller

import (
	"github.com/danielmoisa/auto-runner/internal/driver/keycloak"
	"github.com/danielmoisa/auto-runner/internal/repository"
	"github.com/danielmoisa/auto-runner/internal/utils/cache"
	"github.com/danielmoisa/auto-runner/internal/utils/drive"
	"github.com/danielmoisa/auto-runner/internal/utils/tokenvalidator"
)

type Controller struct {
	Repository            *repository.Repository
	Cache                 *cache.Cache
	Drive                 *drive.Drive
	RequestTokenValidator *tokenvalidator.RequestTokenValidator
	KeycloakClient        *keycloak.Client
	// AttributeGroup        *accesscontrol.AttributeGroup
}

func NewControllerForBackend(repository *repository.Repository, cache *cache.Cache, keycloakClient *keycloak.Client) *Controller { // TODO: attrg *accesscontrol.AttributeGroup, validator *tokenvalidator.RequestTokenValidator
	return &Controller{
		Repository:     repository,
		Cache:          cache,
		KeycloakClient: keycloakClient,
		// Drive:      drive,
		// RequestTokenValidator: validator,
		// AttributeGroup:        attrg,
	}
}

// func NewControllerForBackendInternal(storage *storage.Repository, drive *drive.Drive, validator *tokenvalidator.RequestTokenValidator, attrg *accesscontrol.AttributeGroup) *Controller {
// 	return &Controller{
// 		Repository:               storage,
// 		Drive:                 drive,
// 		RequestTokenValidator: validator,
// 		AttributeGroup:        attrg,
// 	}
// }
