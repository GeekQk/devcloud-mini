package impl_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
)

func TestCreateSecret(t *testing.T) {
	req := &secret.CreateSecretRequest{
		Key:   "api key",
		Value: "api secret",
	}
	res, err := impl.CreateSecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Id)
}

func TestDescribeSecret(t *testing.T) {
	req := &secret.DescribeSecretRequest{
		Id: "cocvuetiika7m4ef8pv0",
	}
	res, err := impl.DescribeSecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	res.Desense()
	t.Log(res)
}

func TestSyncResource(t *testing.T) {
	req := &secret.SyncResourceRequest{
		SecretId: "cod36eh97i69ji730se0",
	}
	err := impl.SyncResource(ctx, req, func(sr *secret.SyncResponse) {
		// 宝塔Linux面板-OhOn[lhins-leelowed]  [OK]
		fmt.Println(sr)
	})
	if err != nil {
		t.Fatal(err)
	}

	// 保证当前进程不退出
	time.Sleep(5 * time.Second)
}
