package secret

import (
	"context"

	resource "github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
)

const (
	AppName = "secret"
)

type Service interface {
	// 录入云商凭证
	CreateSecret(context.Context, *CreateSecretRequest) (*Secret, error)
	// 查询云商凭证
	DescribeSecret(context.Context, *DescribeSecretRequest) (*Secret, error)

	// 使用云商凭证同步资源 Stream
	// SyncResourceHandler 回调参数, 把同步完成的资源, 通过cb交给外部处理
	// 使用 channel   resource chan[SyncResponse]， channel 是锁结构, 回调是无锁结构
	SyncResource(context.Context, *SyncResourceRequest, SyncResourceHandler) error
}

// 通过Hook返回 资源是否同步成功
type SyncResourceHandler func(*SyncResponse)

type SyncResponse struct {
	Id    string
	Name  string
	Error string
}

type SyncResourceRequest struct {
	// secret id
	SecretId string
	// 同步那些区域的资源
	Region []string
	// 同步那些类型的资源
	Resource []resource.TYPE
}

func NewCreateSecretRequest() *CreateSecretRequest {
	return &CreateSecretRequest{
		Regions:       []string{},
		ResourceTypes: []resource.TYPE{},
	}
}
