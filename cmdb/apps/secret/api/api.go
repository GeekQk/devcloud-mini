package api

import (
	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
)

func init() {
	ioc.Api().Registry(&handler{})
}

type handler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl
}

func (h *handler) Name() string {
	return secret.AppName
}

func (h *handler) Init() error {
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
		//  /cmdb/api/v1/secret
		//  cmdb从配置文件中获取的
		//  前缀默认是api 版本默认是v1
		//  secret 对象名称 是在读取的secret.AppName
		//  规则:applicaiton.app.name + "/api/v1/" + secret.AppName,
		ws.GET("/").
			To(h.CreateSecret).
			Doc("创建Secret").
			Metadata("Resource", "Secret").
			Metadata("Action", "Create").
			// 开启了鉴权
			Metadata(label.Auth, label.Enable),
	)

	return nil
}
