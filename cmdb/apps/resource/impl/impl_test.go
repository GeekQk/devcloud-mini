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
	os.Setenv("MONGO_ENDPOINTS", "zdb-cls.cu9pn42.mongodb.net:27017")
	os.Setenv("MONGO_USERNAME", "dbuser")
	os.Setenv("MONGO_PASSWORD", "qiKAI!!395166")
	ioc.DevelopmentSetup()
	fmt.Println(ioc.Controller().List())
	impl = ioc.Controller().Get(resource.AppName).(resource.Service)
}
