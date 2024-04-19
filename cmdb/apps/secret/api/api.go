package api

import (
	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
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
	secret secret.Service
	log    *zerolog.Logger
}

func (h *handler) Name() string {
	return secret.AppName
}

func (h *handler) Init() error {
	h.secret = ioc.Controller().Get(secret.AppName).(secret.Service)
	h.log = log.Sub(h.Name())
	// 对象在初始化的时候 ，就注册给Router
	// Router, 是一个全局对象,
	// 更加当前对象的信息, 来生成一个子路由: r.Group(aa)
	// 关于 GoRestful框架介绍: https://blog.51cto.com/u_15301988/5133427
	// 自动填充对象的 前缀: /app_name/api_prefix/object_version/object_name
	// 在配置文件: app name: cmdb
	// 当前业务模块自动生成的前缀: /cmdb/api/v1/secret
	ws := gorestful.ObjectRouter(h)

	// ws.Route(ws.
	// 	GET("/").
	// 	To(h.MetricHandleFunc).
	// 	Doc("创建Job").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags),
	// )
	// 传递路由规则:  GET  Handler
	ws.Route(
		//  /cmdb/api/v1/secret/xxxx/sync
		ws.POST("/{id}/sync").
			To(h.SyncResource).
			Doc("资源同步").
			// 开启了鉴权
			Metadata(label.Auth, label.Enable).
			// 资源描述
			Metadata(label.Resource, "Secret").
			Metadata(label.Action, "sync").
			// 开启鉴权
			Metadata(label.Permission, label.Enable),
	)

	return nil
}
