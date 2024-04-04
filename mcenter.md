
# 微服务Devcloud研发平台开发: 用户中心-中心化认证
+ 微服务多工程项目组织方式介绍
+ DevCloud需求功能与架构设计
+ 用户中心 认证服务端与客户端开发
+ 用户中心 认证接入中间件开发

详细说明:
```
基于Golang的Workspace组织微服务工程
如何从单体项目中抽象通用模块
基于抽象好的通用模块(二方库)编写一个样例项目
编写helloworld业务模块
注册编写好的业务模块(Controller和API)
程序入口文件, 启动应用程序
补充应用程序配置，包含MongoDB依赖注入
DevCloud需求功能与架构设计
用户中心需求分析
中心化认证
中心化鉴权
中心化认证和鉴权的方案与架构
用户中心技术栈介绍: GoRestful + MongoDB + Grpc
环境准备(主要是MongoDB以及UI工具)
引用用户认证相关所有模块(用户管理，Token颁发等)
本地User/Password 验证流程解读:
Domain的管理
子账号管理(domain下面的账号)
基于本地的令牌的颁发
集成飞书二维码登录验证流程解读:
飞书认证流程解读
准备一个企业账号
创建一个企业自建应用
准备一个飞书登录集成网页
开发飞书登录集成服务端程序
根据返回的code做飞书登录测试
```


## 微服务多工程项目组织方式介绍

+ 放一个代码仓库,    /usercenter /service_b /service_c
+ 一个服务一个仓库,  通过项目组织起来(Git Module)

基于Golang的Workspace组织微服务工程
+ mcenter: 用户中心
+ cmdb: 业务模块, 资源管理
+ audit_log: 操作审计
+ mpaas: 发布中心
+ mperator: k8s 资源同步
+ mflow: 流水线服务

当前已经是一个workspace, 不允许嵌套, 做微服务开发的时候, 自己以项目纬度起workspace

沿用当前 go13 这个workspace


## mcube框架

将单体项目中通用的代码 抽象成 2方库, mcube 微服务开发工具集合

```go
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/mcube/v2/ioc/server"
	"gorm.io/gorm"

	// 开启Health健康检查
	_ "github.com/infraboard/mcube/v2/ioc/apps/health/gin"
	// 开启Metric
	_ "github.com/infraboard/mcube/v2/ioc/apps/metric/gin"
)

func main() {
	// 注册HTTP接口类
	ioc.Api().Registry(&ApiHandler{})

	// 开启配置文件读取配置
	server.DefaultConfig.ConfigFile.Enabled = true
	server.DefaultConfig.ConfigFile.Path = "etc/application.toml"

	// 启动应用
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type ApiHandler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl

	// mysql db依赖
	db *gorm.DB
}

// 覆写对象的名称, 该名称名称会体现在API的路径前缀里面
// 比如: /simple/api/v1/module_a/db_stats
// 其中/simple/api/v1/module_a 就是对象API前缀, 命名规则如下:
// <service_name>/<path_prefix>/<object_version>/<object_name>
func (h *ApiHandler) Name() string {
	return "module_a"
}

// 初始化db属性, 从ioc的配置区域获取共用工具 gorm db对象
func (h *ApiHandler) Init() error {
	h.db = datasource.DB()

	// 进行业务暴露, router 通过ioc
	router := ioc_gin.RootRouter()
	router.GET("/db_stats", h.GetDbStats)
	return nil
}

// 业务功能
func (h *ApiHandler) GetDbStats(ctx *gin.Context) {
	db, _ := h.db.DB()
	ctx.JSON(http.StatusOK, gin.H{
		"data": db.Stats(),
	})
}
```

只关注于业务开发: apps

## 迁移模块到新业务

1. mongodb vs管理插件: Database Client JDBC

我之前开发了很多apps模块, 直接引入到 新业务里面

补充程序配置:
```toml
[mongodb]
  endpoints = ["127.0.0.1:27017"]
  database = "mcentergo13"

[apidoc]
  # Swagger API Doc URL路径, 默认自动生成带前缀的地址比如: /default/api/v1/apidoc 
  # 你也可以在这里直接配置绝对路径 比如: /apidocs.json
  path = "/apidocs.json"
```

启动服务:
```sh
➜  mcenter git:(main) ✗ go run main.go 
2024-03-30T14:25:39+08:00 INFO   cors/gorestful/cors.go:52 > cors enabled component:CORS
2024-03-30T14:25:39+08:00 INFO   apidoc/restful/swagger.go:55 > Get the API Doc using http://127.0.0.1:8080/apidocs.json component:API_DOC
2024-03-30T14:25:39+08:00 INFO   ioc/server/server.go:74 > loaded configs: [app.v1 go_restful_webframework.v1 trace.v1 log.v1 mongodb.v1 go_cache.v1 validator.v1 cache.v1 mcenter.v1 redis.v1 go_restful_cors.v1 grpc.v1 http.v1] component:SERVER
2024-03-30T14:25:39+08:00 INFO   ioc/server/server.go:75 > loaded controllers: [ip2region.v1 counter.v1 domain.v1 endpoint.v1 instances.v1 label.v1 namespace.v1 notify.v1 policy.v1 resource.v1 service.v1 user.v1 oss.v1 role.v1 token.v1] component:SERVER
2024-03-30T14:25:39+08:00 INFO   ioc/server/server.go:76 > loaded apis: [role.v1 label.v1 user/sub.v1 domain.v1 instances.v1 namespace.v1 permission.v1 providers.v1 resource.v1 endpoint.v1 policy.v1 service.v1 code.v1 oauth2.v1 token.v1 account.v1 apidoc.v1] component:SERVER
2024-03-30T14:25:39+08:00 INFO   config/http/http.go:143 > HTTP服务启动成功, 监听地址: 127.0.0.1:8080 component:HTTP
2024-03-30T14:25:39+08:00 INFO   config/grpc/grpc.go:138 > GRPC 服务监听地址: 127.0.0.1:18080 component:GRPC
2024-03-30T14:35:10+08:00 ERROR  token/security/checker.go:195 > get user error, user describe_by:USER_NAME  username:"admin" not found, use default setting to check component:"LOGIN SECURITY"
2024-03-30T14:35:10+08:00 DEBUG  token/security/checker.go:51 > max failed retry lock check enabled, checking ... component:"LOGIN SECURITY"
2024-03-30T14:35:10+08:00 ERROR  token/security/checker.go:56 > get key abnormal_admin from cache error, Key not found. component:"LOGIN SECURITY"
2024-03-30T14:35:10+08:00 DEBUG  token/security/checker.go:60 > retry times: 0, retry limite: 5 component:"LOGIN SECURITY"
2024-03-30T14:35:10+08:00 ERROR  token/security/checker.go:195 > get user error, user describe_by:USER_NAME  username:"admin" not found, use default setting to check component:"LOGIN SECURITY"
2024-03-30T14:35:10+08:00 DEBUG  token/security/checker.go:182 > ip limite check disabled, don't check component:"LOGIN SECURITY"
2024-03-30T14:35:10+08:00 DEBUG  token/impl/token.go:105 > security check complete component:TOKEN
```


2. 尝试用  swagger ui 进行 API接口文档展示: http://127.0.0.1:8080/apidocs.json

```sh
docker pull swaggerapi/swagger-ui
docker run -p 80:8080 swaggerapi/swagger-ui

# 后端开启cros支持:  _ "github.com/infraboard/mcube/v2/ioc/config/cors/gorestful"	
# 访问本地 http://127.0.0.1 访问Swagger UI
```

## 初始化admin账号

引入cmd工具
```go
	// 加载所有模块
	_ "gitlab.com/go-course-project/go13/devcloud-mini/mcenter/apps"
	// 开启API Doc
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
	// 支持跨越
	_ "github.com/infraboard/mcube/v2/ioc/config/cors/gorestful"

	// 注入初始化命令:
	"github.com/infraboard/mcenter/cmd/initial"
```

初始化管理员账号:
```sh
➜  mcenter git:(main) ✗ go run main.go init          
2024-03-30T14:53:36+08:00 INFO   cors/gorestful/cors.go:52 > cors enabled component:CORS
2024-03-30T14:53:36+08:00 INFO   apidoc/restful/swagger.go:55 > Get the API Doc using http://127.0.0.1:8080/apidocs.json component:API_DOC
? 请输入公司(组织)名称: 基础设施服务中心
? 请输入管理员用户名称: admin
? 请输入管理员密码: ******
? 再次输入管理员密码: ******
初始化域:           default [成功]
初始化系统管理员:     admin [成功]
初始化空间:         default [成功]
初始化空间:          system [成功]
2024-03-30T14:54:06+08:00 DEBUG  role/impl/permission.go:19 > query permission filter: map[role_id:15ed3552927dd711] component:ROLE
2024-03-30T14:54:06+08:00 DEBUG  role/impl/permission.go:19 > query permission filter: map[role_id:15ed3552927dd711] component:ROLE
2024-03-30T14:54:06+08:00 DEBUG  role/impl/permission.go:19 > query permission filter: map[role_id:15ed3552927dd711] component:ROLE
初始化角色:           admin [成功]
2024-03-30T14:54:06+08:00 DEBUG  role/impl/permission.go:19 > query permission filter: map[role_id:7055db386b702304] component:ROLE
2024-03-30T14:54:06+08:00 DEBUG  role/impl/permission.go:19 > query permission filter: map[role_id:7055db386b702304] component:ROLE
2024-03-30T14:54:06+08:00 DEBUG  role/impl/permission.go:19 > query permission filter: map[role_id:7055db386b702304] component:ROLE
初始化角色:         visitor [成功]
初始化服务:          maudit [成功]
初始化服务:           mflow [成功]
初始化服务:       moperator [成功]
初始化服务:           mpaas [成功]
初始化服务:            cmdb [成功]
初始化标签:             Env [成功]初始化标签:   ResourceGroup [成功]初始化标签:       UserGroup [成功]Error: 初始化标签失败: inserted a label document error, write exception: write errors: [E11000 duplicate key error collection: mcentergo13.label index: _id_ dup key: { _id: "9f29507fbd437645" }]
```

## 登录流程解读

[](./docs/login.drawio)


## 用户校验

开发一个中间件, 中间件 校验Token的逻辑是 使用GRPC


1. 用户中心 用户校验: 用户中心自己的安全才有保证, 进程内调用 内部token 来进行 Token的校验

mcenter 服务端 认证中间件
```go
// 查看auth 校验逻辑: https://github.com/infraboard/mcenter/blob/master/middlewares/http/auth.go
// 查看auth 校验逻辑: https://github.com/infraboard/mcenter/blob/master/middlewares/grpc/auth.go
import (
	// http  认证拦截器
	_ "github.com/infraboard/mcenter/middlewares/grpc"
	// grpc 认证拦截器
	_ "github.com/infraboard/mcenter/middlewares/http"
)
```

2. 其他服务 接入用户中心 用户校验, 进程间调用, 通过grpc 调用外部token 来进行校验

第三方服务如何接入 中心:  外部Token认证中间件
```go
import (
	// 基于Grpc Token服务 Validate Token的认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/gorestful"
	// 基于Grpc Service服务的 Validate ClientId/ClientSecret 的grpc认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/grpc"
)
```

## cmdb 接入mcenter 进行鉴权


cmdb Handler
```go
// 写一个 go restful 框架的Handler
// 使用的v3版本
//
//	Gin （ctx *gin.Context)
//
// GoRestful (r *restful.Request, w *restful.Response)
func (h *handler) CreateSecret(r *restful.Request, w *restful.Response) {
	// 1. 获取请求 Request Entity
	// 2. Gin Bind
	// r.ReadEntity()

	// 2. 处理请求

	// 3. 返回请求
	// c.JSON()
	// w.WriteAsJson()
	// w.WriteError()
	w.WriteAsJson(map[string]any{"code": 0})
}
```

路由注册
```go
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
		ws.GET("/").
			To(h.CreateSecret).
			Doc("创建Secret").
			Metadata("Resource", "Secret").
			Metadata("Action", "Create"),
	)

	return nil
}
```

```sh
➜  go13 git:(main) ✗ curl http://127.0.0.1:7010/cmdb/api/v1/secret
{
 "code": 0
}%    
```

补上认证中间件
```go
import (
	// 基于Grpc Token服务 Validate Token的认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/gorestful"
	// 基于Grpc Service服务的 Validate ClientId/ClientSecret 的grpc认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/grpc"
)
```

接口开启开启, 需要通过路由装饰: 添加Meta信息, 对接口做备注: Metadata(label.Auth, label.Enable)
```go
	ws.Route(
		//  /cmdb/api/v1/secret
		ws.GET("/").
			To(h.CreateSecret).
			Doc("创建Secret").
			Metadata("Resource", "Secret").
			Metadata("Action", "Create").
			// 开启了鉴权
			Metadata(label.Auth, label.Enable),
	)
```

在此尝试携带token进行认证

```
[mcenter]
address = "127.0.0.1:17080"
client_id = "M7LmSqFgitbqE9sckoZP2vZY"
client_secret = "ZcpbNEftiIP04ErVdEn29mHosTYYr0FL"
```





