package impl_test

import (
	"testing"
	"time"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/infraboard/mcenter/apps/domain"
	"github.com/infraboard/mcenter/apps/namespace"
)

func TestSaveResource(t *testing.T) {
	ins := resource.NewResource()
	ins.Meta.Id = "test01"
	ins.Meta.Domain = domain.DEFAULT_DOMAIN
	ins.Meta.Namespace = namespace.DEFAULT_NAMESPACE
	ins.Meta.SyncAt = time.Now().Unix()
	ins.Spec.Name = "cvm01"
	ins.Spec.Type = resource.TYPE_HOST.String()
	res, err := impl.Save(ctx, ins)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestSearchResource(t *testing.T) {
	req := resource.NewSearchRequest()
	req.Keywords = "cvm"
	res, err := impl.Search(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
