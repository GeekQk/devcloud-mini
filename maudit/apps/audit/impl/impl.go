package impl

import (
	"github.com/GeekQk/devcloud-mini/maudit/apps/audit"
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

	col *mongo.Collection
	log *zerolog.Logger
}

func (i *impl) Name() string {
	return audit.AppName
}

func (i *impl) Init() error {
	// resource表
	i.col = ioc_mongo.DB().Collection("audit_log")
	// 模块日志
	i.log = log.Sub(i.Name())
	return nil
}
