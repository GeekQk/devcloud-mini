package main

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/server"

	// 引入所有的业务模块
	_ "github.com/GeekQk/devcloud-mini/cmdb/apps"

	// 开启API Doc
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
	// 支持跨越
	_ "github.com/infraboard/mcube/v2/ioc/config/cors/gorestful"
	// // 基于Grpc Token服务 Validate Token的认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/gorestful"
	// 基于Grpc Service服务的 Validate ClientId/ClientSecret 的grpc认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/grpc"
)

func main() {
	// 开启配置文件读取配置
	server.DefaultConfig.ConfigFile.Enabled = true
	server.DefaultConfig.ConfigFile.Path = "etc/application.toml"

	// 启动应用
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}

}
