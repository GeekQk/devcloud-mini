package impl_test

import (
	"testing"

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
