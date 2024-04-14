package tecent

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/GeekQk/devcloud-mini/cmdb/provider"
)

// 提供CVM同步 就是服务器同步
type CvmProvider struct {
	provider.ResourceSyncConfig
}

// 实现资源同步 Stream
// 1. 拉去资源
// 2. 转化
// 3. 处理
func (c *CvmProvider) Sync(ctx context.Context, hanleFunc provider.SyncResourceHandler) error {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		c.ApiKey,
		c.ApiSecret,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	// 推荐使用北极星，相关指引可访问如下链接
	// https://git.woa.com/tencentcloud-internal/tencentcloud-sdk-go#%E5%8C%97%E6%9E%81%E6%98%9F
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := lighthouse.NewClient(credential, c.Region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := lighthouse.NewDescribeInstancesRequest()

	// 返回的resp是一个DescribeInstancesResponse的实例，与请求对象对应
	response, err := client.DescribeInstances(request)
	if err != nil {
		return err
	}

	//  转换
	for _, ins := range response.Response.InstanceSet {
		res, err := c.TransferInstance(ins)
		if err != nil {
			log.L().Error().Msgf("transfer instance error, %s", err)
			continue
		}
		hanleFunc(ctx, res)
	}
	return nil
}

// 统一到Resource模型
func (c *CvmProvider) TransferInstance(
	ins *lighthouse.Instance) (
	*resource.Resource, error) {
	res := resource.NewResource()

	res.Meta.Id = *ins.InstanceId
	res.Meta.SyncAt = time.Now().Unix()

	// Spec
	res.Spec.Vendor = resource.VENDOR_TENCENT
	res.Spec.ResourceType = resource.TYPE_HOST
	res.Spec.Region = c.Region
	res.Spec.Zone = *ins.Zone
	res.Spec.Name = *ins.InstanceName
	res.Spec.Category = *ins.BundleId
	res.Spec.Type = *ins.BundleId

	// 处理时间: ExpiredTime 2024-12-22T08:24:41Z
	et, err := time.Parse("2006-01-02T15:04:05Z", *ins.ExpiredTime)
	if err != nil {
		return nil, err
	}
	res.Spec.ExpireAt = et.Unix()
	res.Spec.Cpu = int32(*ins.CPU)
	res.Spec.Memory = int32(*ins.Memory) * 1024
	res.Spec.Storage = int32(*ins.SystemDisk.DiskSize)
	res.Spec.BandWidth = int32(*ins.InternetAccessible.InternetMaxBandwidthOut)

	// "OsName": "CentOS 7.9 64bit",
	// "Platform": "CENTOS",
	// "PlatformType": "LINUX_UNIX",
	// 主机特有的属性
	// extra.os_platform = 'CENTOS'
	res.Spec.Extra["os_name"] = *ins.OsName
	res.Spec.Extra["os_platform"] = *ins.Platform
	res.Spec.Extra["os_platform_type"] = *ins.PlatformType

	// Status, 统一多个厂商的状态
	res.Status.Phase = *ins.InstanceState
	// 内网IP
	for _, pa := range ins.PrivateAddresses {
		res.Status.PrivateAddress = append(res.Status.PrivateAddress, *pa)
	}
	//公网IP
	for _, pa := range ins.PublicAddresses {
		res.Status.PublicAddress = append(res.Status.PublicAddress, *pa)
	}
	return res, nil
}
