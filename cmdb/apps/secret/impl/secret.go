package impl

import (
	"context"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

// 录入云商凭证
func (i *impl) CreateSecret(
	ctx context.Context,
	in *secret.CreateSecretRequest) (
	*secret.Secret, error) {
	// 1. 校验请求
	if err := validator.Validate(in); err != nil {
		return nil, err
	}

	// 2. 构建实例, tk 获取
	ins := &secret.Secret{
		Id:   xid.New().String(),
		Spec: in,
	}

	// 3. 存储
	if err := ins.Encrypt(); err != nil {
		return nil, err
	}
	if _, err := i.col.InsertOne(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

// 查询云商凭证
func (i *impl) DescribeSecret(
	ctx context.Context,
	in *secret.DescribeSecretRequest) (
	*secret.Secret, error) {

	ins := &secret.Secret{}
	if err := i.col.FindOne(ctx, bson.M{"_id": in.Id}).Decode(ins); err != nil {
		return nil, err
	}
	if err := ins.Decrypt(); err != nil {
		return nil, err
	}
	return ins, nil
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
