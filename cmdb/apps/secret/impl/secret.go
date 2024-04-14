package impl

import (
	"context"
	"time"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
	"github.com/GeekQk/devcloud-mini/cmdb/provider"
	"github.com/GeekQk/devcloud-mini/cmdb/provider/tecent"
	"github.com/infraboard/mcube/v2/ioc/config/log"
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

	// 1. 找到secret
	req := &secret.DescribeSecretRequest{Id: in.SecretId}
	secretIns, err := i.DescribeSecret(ctx, req)
	if err != nil {
		return err
	}

	// 2. 初始化CVMProvider, 资源提供商列表
	resourceProviders := []provider.ResourceProvider{}
	//默认使用secret.Spec.Regions
	region := []string{}
	if len(in.Region) == 0 {
		region = secretIns.Spec.Regions
	}
	for _, region := range region {
		conf := provider.ResourceSyncConfig{
			ApiKey:    secretIns.Spec.Key,
			ApiSecret: secretIns.Spec.Value,
			Region:    region,
		}
		resourceProviders = append(resourceProviders,
			&tecent.CvmProvider{ResourceSyncConfig: conf})
	}

	// 3. 使用 CVMProvider 进行资源同步
	// 同步 进行资源同步，还是异步进行资源同步 (task)
	// ctx 请求的context, 任务还需要继续进行
	gctx, _ := context.WithTimeout(context.Background(), 1*time.Hour)
	for _, rp := range resourceProviders {
		// //异步请求:第一种方式
		// rp.Sync(gctx, i.SyncResoruce)

		//同步请求:第二种方式
		rp.Sync(gctx, func(ctx context.Context, r *resource.Resource) {
			_, err := i.resource.Save(ctx, r)
			if err != nil {
				//同步失败
				log.L().Error().Msgf("save resource error, %s", err)
			} else {
				//同步成功
				rh(&secret.SyncResponse{
					Id:   r.Meta.Id,
					Name: r.Spec.Name,
				})
			}
		})
	}
	return nil
}

// 4. 将同步的结果 调用resource 录入
func (i *impl) SyncResoruce(ctx context.Context, r *resource.Resource) {
	_, err := i.resource.Save(ctx, r)
	if err != nil {
		log.L().Error().Msgf("save resource error, %s", err)
	}
}
