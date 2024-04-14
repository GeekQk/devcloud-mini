package impl

import (
	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	// 需要托管Ioc
	ioc.ObjectImpl
	// 基础GRPC Server实现
	resource.UnimplementedRPCServer

	col      *mongo.Collection
	log      *zerolog.Logger
	resource resource.Service
}

func (i *impl) Name() string {
	return secret.AppName
}

func (i *impl) Init() error {
	// resource表
	i.col = ioc_mongo.DB().Collection("secret")
	// 模块日志
	i.log = log.Sub(i.Name())
	// 依赖的资源服务
	i.resource = ioc.Controller().Get(resource.AppName).(resource.Service)
	return nil
}
