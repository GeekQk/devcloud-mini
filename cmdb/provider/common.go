package provider

import (
	"context"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
)

type ResourceSyncConfig struct {
	ApiKey    string
	ApiSecret string
	Region    string
}

type SyncResourceHandler func(context.Context, *resource.Resource)

type ResourceProvider interface {
	Sync(ctx context.Context, hanleFunc SyncResourceHandler) error
}
