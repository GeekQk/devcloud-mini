package main

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/server"

	// 引入业务模块
	_ "github.com/GeekQk/devcloud-mini/maudit/apps"
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
