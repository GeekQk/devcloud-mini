package main

import (
	"context"

	//cmd
	"github.com/spf13/cobra"

	//注入初始化命令
	"github.com/infraboard/mcenter/cmd/initial"

	//引入server
	"github.com/infraboard/mcube/v2/ioc/server"
	"github.com/infraboard/mcube/v2/ioc/server/cmd"

	//引入api doc
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"

	//加载init方法的所有模块  ->apps->mcube其他模块
	_ "github.com/GeekQk/devcloud-mini/mcenter/apps"

	//支持跨域
	_ "github.com/infraboard/mcube/v2/ioc/config/cors/gorestful"

	// http  认证拦截器
	_ "github.com/infraboard/mcenter/middlewares/grpc"
	// grpc 认证拦截器
	_ "github.com/infraboard/mcenter/middlewares/http"
)

func main() {
	// 开启配置文件读取配置
	// server.DefaultConfig.ConfigFile.Enabled = true
	// server.DefaultConfig.ConfigFile.Path = "etc/application.toml"
	// // 启动应用
	// err := server.Run(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// // 全局Root CMD
	cmd.Root.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "example API服务",
			Run: func(cmd *cobra.Command, args []string) {
				// // 开启配置文件读取配置
				server.DefaultConfig.ConfigFile.Enabled = true
				server.DefaultConfig.ConfigFile.Path = "etc/application.toml"
				cobra.CheckErr(server.Run(context.Background()))
			},
		},
	)

	cmd.Root.AddCommand(initial.Cmd)

	// 启动
	cmd.Execute()
}
