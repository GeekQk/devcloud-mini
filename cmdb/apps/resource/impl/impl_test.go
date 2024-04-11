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
	os.Setenv("MONGO_DATABASE", "cmdb")
	ioc.DevelopmentSetup()

	fmt.Println(ioc.Controller().List())
	impl = ioc.Controller().Get(resource.AppName).(resource.Service)
}
