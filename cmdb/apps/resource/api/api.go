package api

import (
	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&handler{})
}

type handler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl

	resource resource.Service
	log      *zerolog.Logger
}

func (h *handler) Name() string {
	return resource.AppName
}

func (h *handler) Init() error {
	h.resource = ioc.Controller().Get(resource.AppName).(resource.Service)
	h.log = log.Sub(h.Name())

	ws := gorestful.ObjectRouter(h)

	ws.Route(
		//  /cmdb/api/v1/resource
		ws.GET("/").
			To(h.SearchResource).
			Doc("资源搜索").
			// 开启了鉴权
			Metadata(label.Auth, label.Enable).
			// 资源描述
			Metadata(label.Resource, "Resource").
			Metadata(label.Action, "list").
			// 开启鉴权
			Metadata(label.Permission, label.Disable),
	)
	return nil
}
