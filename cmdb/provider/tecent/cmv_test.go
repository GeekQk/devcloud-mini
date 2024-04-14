package tecent_test

import (
	"context"
	"testing"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/GeekQk/devcloud-mini/cmdb/provider"
	"github.com/GeekQk/devcloud-mini/cmdb/provider/tecent"
)

var ctx = context.Background()

// var client = &tecent.CvmProvider{
// 	provider.ResourceSyncConfig{
// 		ApiKey:    os.Getenv("API_KEY"),
// 		ApiSecret: os.Getenv("API_SECRET"),
// 		Region:    "ap-guangzhou",
// 	},
// }

var client = &tecent.CvmProvider{
	provider.ResourceSyncConfig{
		ApiSecret: "zthUCjCoZHyjRBykVCj2CCZw6UtAF6sr",
		ApiKey:    "AKIDDSm4OweGzytdFetBNXaRMPVgvaRhqqRI",
		Region:    "ap-guangzhou",
	},
}

func TestSync(t *testing.T) {
	err := client.Sync(ctx, func(ctx context.Context, r *resource.Resource) {
		t.Log(r)
	})
	if err != nil {
		t.Fatal(err)
	}
}
