package impl_test

import (
	"testing"

	"gitlab.com/go-course-project/go13/devcloud-mini/cmdb/apps/secret"
)

func TestCreateSecret(t *testing.T) {
	req := &secret.CreateSecretRequest{}
	res, err := impl.CreateSecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
