package impl_test

import (
	"context"
	"fmt"
	"os"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/infraboard/mcube/v2/ioc"

	// 加载所有模块
	_ "github.com/GeekQk/devcloud-mini/cmdb/apps"
)

var (
	impl resource.Service
	ctx  = context.Background()
)

func init() {
	// 开启配置文件读取配置
	os.Setenv("MONGO_DATABASE", "cmdb")
	os.Setenv("MONGO_ENDPOINTS", "dds-bp1ef8a2abb33ef41762-pub.mongodb.rds.aliyuncs.com:3717")
	os.Setenv("MONGO_USERNAME", "dbuser")
	os.Setenv("MONGO_PASSWORD", "qiKAI!!395166")
	ioc.DevelopmentSetup()
	fmt.Println(ioc.Controller().List())
	impl = ioc.Controller().Get(resource.AppName).(resource.Service)
}
