package controller

import (
	"github.com/danielmoisa/matiq/internal/driver/keycloak"
	"github.com/danielmoisa/matiq/internal/repository"
	"github.com/danielmoisa/matiq/internal/utils/cache"
	"github.com/danielmoisa/matiq/internal/utils/drive"
	"github.com/danielmoisa/matiq/internal/utils/tokenvalidator"
)

type Controller struct {
	Repository            *repository.Repository
	Cache                 *cache.Cache
	Drive                 *drive.Drive
	RequestTokenValidator *tokenvalidator.RequestTokenValidator
	KeycloakClient        *keycloak.Client
}

func NewControllerForBackend(repository *repository.Repository, cache *cache.Cache, keycloakClient *keycloak.Client) *Controller { // TODO: attrg *accesscontrol.AttributeGroup, validator *tokenvalidator.RequestTokenValidator
	return &Controller{
		Repository:     repository,
		Cache:          cache,
		KeycloakClient: keycloakClient,
	}
}
