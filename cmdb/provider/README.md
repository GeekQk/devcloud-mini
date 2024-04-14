# 资源供应商


[](./sync.drawio)


```go
// 提供CVM同步
type CvmProvider struct {
}

// 实现资源同步 Stream
func (c *CvmProvider) Sync(ctx context.Context) error {
	return nil
}

```

API 使用与调试: https://console.cloud.tencent.com/api/explorer?Product=cvm&Version=2017-03-12&Action=DescribeRegions


```go
package main

import (
        "fmt"

        "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
        "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
        "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
        lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

func main() {
        // 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
        // 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
        // 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
        credential := common.NewCredential(
                "SecretId",
                "SecretKey",
        )
        // 实例化一个client选项，可选的，没有特殊需求可以跳过
        cpf := profile.NewClientProfile()
        // 推荐使用北极星，相关指引可访问如下链接
        // https://git.woa.com/tencentcloud-internal/tencentcloud-sdk-go#%E5%8C%97%E6%9E%81%E6%98%9F
        cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
        // 实例化要请求产品的client对象,clientProfile是可选的
        client, _ := lighthouse.NewClient(credential, "ap-guangzhou", cpf)

        // 实例化一个请求对象,每个接口都会对应一个request对象
        request := lighthouse.NewDescribeInstancesRequest()
        

        // 返回的resp是一个DescribeInstancesResponse的实例，与请求对象对应
        response, err := client.DescribeInstances(request)
        if _, ok := err.(*errors.TencentCloudSDKError); ok {
                fmt.Printf("An API error has returned: %s", err)
                return
        }
        if err != nil {
                panic(err)
        }
        // 输出json格式的字符串回包
        fmt.Printf("%s", response.ToJsonString())
} 
```