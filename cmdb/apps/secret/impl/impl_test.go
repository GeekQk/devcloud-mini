package impl_test

import (
	"context"
	"fmt"
	"os"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
	"github.com/infraboard/mcube/v2/ioc"
)

var (
	impl secret.Service
	ctx  = context.Background()
)

func init() {
	os.Setenv("MONGO_DATABASE", "cmdb")
	ioc.DevelopmentSetup()

	fmt.Println(ioc.Controller().List())
	impl = ioc.Controller().Get(secret.AppName).(secret.Service)
}
