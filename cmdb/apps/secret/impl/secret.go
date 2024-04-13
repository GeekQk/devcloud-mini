package impl

import (
	"context"

	"gitlab.com/go-course-project/go13/devcloud-mini/cmdb/apps/secret"
)

// 录入云商凭证
func (i *impl) CreateSecret(
	ctx context.Context,
	in *secret.CreateSecretRequest) (
	*secret.Secret, error) {
	return nil, nil
}

// 查询云商凭证
func (i *impl) DescribeSecret(
	ctx context.Context,
	in *secret.DescribeSecretRequest) (
	*secret.Secret, error) {
	return nil, nil
}

// 使用云商凭证同步资源 Stream
// SyncResourceHandler 回调参数, 把同步完成的资源, 通过cb交给外部处理
// 使用 channel   resource chan[SyncResponse]， channel 是锁结构, 回调是无锁结构
func (i *impl) SyncResource(
	ctx context.Context,
	in *secret.SyncResourceRequest,
	rh secret.SyncResourceHandler) error {
	return nil
}
