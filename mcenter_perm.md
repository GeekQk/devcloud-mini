# 中心化鉴权

```
基于中间件的接入架构设计
基于用户中心GRPC客户端开发认证中间件
初始化一个样例服务(cmdb) 编写一个secret管理模块, 用于验证接入测试
测试cmdb服务能否使用认证中间件接入用户中心进行统一认证
RBAC与PBAC的鉴权模型解读
基于PBAC鉴权模型设计以及鉴权期望
引用用户鉴权相关所有模块
用户鉴权流程解读:
接入中心化鉴权的服务(cmdb)注册功能列表
创建权限测试角色 管理员/访客
创建2个测试用户 admin01/visotor01
为2个测试用户 添加授权策略
验证鉴权功能: 分别校验admin01/visotor01用户是否有创建secret的权限
利用RPC客户端 封装鉴权逻辑, 中间件补充鉴权
更新最新的接入中间件, 测试cmdb 接入用户中心后的鉴权效果
```

[](./arch.drawio)


## 功能注册

1. mcenter endpoint 管理模块开发: 功能注册RPC: RegistryEndpoint

```proto
// RPC endpoint管理
service RPC {
	rpc DescribeEndpoint(DescribeEndpointRequest) returns(Endpoint);
	rpc QueryEndpoints(QueryEndpointRequest) returns(EndpointSet);
	rpc RegistryEndpoint(RegistryRequest) returns(RegistryResponse);
}
```

2. 客户端注册功能列表: 

功能注册中间件:
```go
import (
    _ 	"github.com/infraboard/mcenter/clients/rpc/middleware/registry/endpoint"
)
```

遍历 Container所有web service 转化所有的Route ---> Entry, 注册这个服务所有的 Entry
```go
// gorestful 的跟理由传递给过滤
func (r *EndpointRegister) Registry(ctx context.Context, c *restful.Container, version string) error {
	entries := []*endpoint.Entry{}
	wss := c.RegisteredWebServices()
	for i := range wss {
		es := endpoint.TransferRoutesToEntry(wss[i].Routes())
		entries = append(entries, endpoint.GetPRBACEntry(es)...)
	}

	req := endpoint.NewRegistryRequest(version, entries)
	_, err := r.c.Endpoint().RegistryEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}
```

```go
func NewEntryFromRestRoute(route restful.Route) *Entry {
	entry := NewDefaultEntry()
	entry.FunctionName = route.Operation
	entry.Method = route.Method
	entry.LoadMeta(route.Metadata)
	entry.Path = route.Path

	entry.Path = entry.UniquePath()
	return entry
}
```

endpoints 模块数据库里面的 功能条目: 
```json
{
    "_id": "bd12ef56c19e085d",
    "create_at": 1712371188,
    "update_at": 1712371188,
    "service_id": "3f15485049fb5f43",
    "version": "[ ]",
    "function_name": "CreateSecret",
    "path": "GET./cmdb/api/v1/secret/",
    "method": "GET",
    "resource": "Secret",
    "auth_enable": "true",
    "code_enable": "false",
    "permission_mode": "PRBAC",
    "permission_enable": "false",
    "allow": null,
    "audit_log": "false",
    "required_namespace": "false",
    "labels": "{\"action\":\"create\"}",
    "extension": "{\"domain\":\"default\",\"namespace\":\"system\"}"
}
```

## 用户授权

1. 创建用户

```sh
# admin01
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/user/sub' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=4mifKdQnf5BUzkPBmNr15x5r' \
--data '{
    "username": "admin01",
    "password": "123456"
}'

# visitor01
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/user/sub' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=4mifKdQnf5BUzkPBmNr15x5r' \
--data '{
    "username": "visitor01",
    "password": "123456"
}'
```

2. 创建2个角色

```go
type PermissionSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 创建人
	// @gotags: bson:"create_by" json:"create_by"
	CreateBy string `protobuf:"bytes,1,opt,name=create_by,json=createBy,proto3" json:"create_by" bson:"create_by"`
	// 权限描述
	// @gotags: bson:"desc" json:"desc"
	Desc string `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc" bson:"desc"`
	// 效力, Allow/Deny
	// @gotags: bson:"effect" json:"effect"
	Effect EffectType `protobuf:"varint,4,opt,name=effect,proto3,enum=infraboard.mcenter.role.EffectType" json:"effect" bson:"effect"`
	// 服务ID
	// @gotags: bson:"service_id" json:"service_id"
	ServiceId string `protobuf:"bytes,5,opt,name=service_id,json=serviceId,proto3" json:"service_id" bson:"service_id"`
	// 资源列表
	// @gotags: bson:"resource_name" json:"resource_name"
	ResourceName string `protobuf:"bytes,6,opt,name=resource_name,json=resourceName,proto3" json:"resource_name" bson:"resource_name"`
	// 维度
	// @gotags: bson:"label_key" json:"label_key"
	LabelKey string `protobuf:"bytes,7,opt,name=label_key,json=labelKey,proto3" json:"label_key" bson:"label_key"`
	// 适配所有值
	// @gotags: bson:"match_all" json:"match_all"
	MatchAll bool `protobuf:"varint,8,opt,name=match_all,json=matchAll,proto3" json:"match_all" bson:"match_all"`
	// 标识值
	// @gotags: bson:"label_values" json:"label_values"
	LabelValues []string `protobuf:"bytes,9,rep,name=label_values,json=labelValues,proto3" json:"label_values" bson:"label_values"`
}
```

+ Visitor: 访客, 对所有资源可读权限,  Allow  所有服务: * 所有的资源: * LabelKey: action, LabelValues: ["get", "list"]
```sh
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/role' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=4mifKdQnf5BUzkPBmNr15x5r' \
--data '{
    "name": "visitor01",
    "permissions": [
        {
            "service_id": "*",
            "resource_name": "*",
            "label_key": "action",
            "label_values": ["get", "list"]
        }
    ]
}'
```

+ admin01: 管理员, 所有资源的所有权限,  Allow  所有服务: * 所有的资源: * LabelKey: *, MatchALl

```sh
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/role' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=4mifKdQnf5BUzkPBmNr15x5r' \
--data '{
    "name": "admin01",
    "permissions": [
        {
            "service_id": "*",
            "resource_name": "*",
            "label_key": "*",
            "match_all": true
        }
    ]
}'
```

3. 给用户授权(添加授权策略)

```go
// CreatePolicyRequest 创建策略的请求
type CreatePolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields


	// 用户Id
	// @gotags: bson:"user_id" json:"user_id" validate:"required,lte=120"
	UserId string `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id" bson:"user_id" validate:"required,lte=120"`
	// 范围
	// @gotags: bson:"namespace" json:"namespace" validate:"lte=120"
	Namespace string `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace" bson:"namespace" validate:"lte=120"`
	// 角色Id
	// @gotags: bson:"role_id" json:"role_id" validate:"required,lte=40"
	RoleId string `protobuf:"bytes,5,opt,name=role_id,json=roleId,proto3" json:"role_id" bson:"role_id" validate:"required,lte=40"`

	// 创建者
	// @gotags: bson:"create_by" json:"create_by"
	CreateBy string `protobuf:"bytes,1,opt,name=create_by,json=createBy,proto3" json:"create_by" bson:"create_by"`
	// 策略所属域
	// @gotags: bson:"domain" json:"domain"
	Domain string `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain" bson:"domain"`
	// 该角色的生效范围
	// @gotags: bson:"scope" json:"scope"
	Scope []*resource.LabelRequirement `protobuf:"bytes,6,rep,name=scope,proto3" json:"scope" bson:"scope"`
	// 策略过期时间
	// @gotags: bson:"expired_time" json:"expired_time"
	ExpiredTime int64 `protobuf:"varint,7,opt,name=expired_time,json=expiredTime,proto3" json:"expired_time" bson:"expired_time"`
	// 只读策略, 不允许用户修改, 一般用于系统管理
	// @gotags: bson:"read_only" json:"read_only"
	ReadOnly bool `protobuf:"varint,8,opt,name=read_only,json=readOnly,proto3" json:"read_only" bson:"read_only"`
	// 启用该策略
	// @gotags: bson:"enabled" json:"enabled"
	Enabled bool `protobuf:"varint,9,opt,name=enabled,proto3" json:"enabled" bson:"enabled"`
	// 扩展属性
	// @gotags: bson:"extra" json:"extra"
	Extra map[string]string `protobuf:"bytes,14,rep,name=extra,proto3" json:"extra" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"extra"`
	// 标签
	// @gotags: bson:"labels" json:"labels"
	Labels map[string]string `protobuf:"bytes,15,rep,name=labels,proto3" json:"labels" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"labels"`
}
```


3.1. admin用户授权
```sh
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/policy' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=4mifKdQnf5BUzkPBmNr15x5r' \
--data-raw '{
    "user_id": "admin01@default",
    "namespace": "default",
    "role_id": "9c38a5db4f9480fe"
}'
```

3.2 visitor 用户授权
```sh
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/policy' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=4mifKdQnf5BUzkPBmNr15x5r' \
--data-raw '{
    "user_id": "visitor01@default",
    "namespace": "default",
    "role_id": "aea82c6cd4ce35cf"
}'
```

## 鉴权验证

1. 开启鉴权: Metadata(label.Permission, label.Enable)
```go
	ws.Route(
		//  /cmdb/api/v1/secret
		ws.GET("/").
			To(h.CreateSecret).
			Doc("创建Secret").
			Metadata("Resource", "Secret").
			Metadata("Action", "Create").
			// 开启了鉴权
			Metadata(label.Auth, label.Enable).
			// 资源描述
			Metadata(label.Resource, "Secret").
			Metadata(label.Action, label.Create.Value()).
			// 开启鉴权
			Metadata(label.Permission, label.Enable),
	)
```

Secret资源: admin01能访问, visitor01是无法访问的

2.1 admin01能访问
```sh
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/token' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=DhRcVREClCdoqU4N03GMtKbr' \
--data '{
  "username": "admin01",
  "password": "123456"
}'

curl --location 'http://127.0.0.1:7010/cmdb/api/v1/secret' \
--header 'Cookie: mcenter.access_token=DhRcVREClCdoqU4N03GMtKbr'
```

2.2 visitor01是无法访问的

```sh
curl --location 'http://127.0.0.1:7080/mcenter/api/v1/token' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=DhRcVREClCdoqU4N03GMtKbr' \
--data '{
  "username": "visitor01",
  "password": "123456"
}'

curl --location 'http://127.0.0.1:7010/cmdb/api/v1/secret' \
--header 'Cookie: mcenter.access_token=DhRcVREClCdoqU4N03GMtKbr'
```

```json
{
    "namespace": "cmdb",
    "http_code": 403,
    "error_code": 403,
    "reason": "访问未授权",
    "message": "in namespace default, role [visitor01@default] has no permission access endpoint: GET./cmdb/api/v1/secret/",
    "meta": null,
    "data": null
}
```

## 鉴权流程

[](./docs/perm.drawio)


## 总结

mcenter  微服务架构下, 中心化的认证和鉴权, 其他子服务 是如何接入 微服务的 用户中心

通过导入中间件接入用户中心:
```go
	// 基于Grpc Token服务 Validate Token的认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/gorestful"

	// 基于Grpc Service服务的 Validate ClientId/ClientSecret 的grpc认证中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/auth/grpc"

	// 服务功能列表注册中间件
	_ "github.com/infraboard/mcenter/clients/rpc/middleware/registry/endpoint"
```

通过路由装饰 来开启认证和鉴权, 路由描述
```go
	ws.Route(
		//  /cmdb/api/v1/secret
		ws.GET("/").
			To(h.CreateSecret).
			Doc("创建Secret").
			// 开启了鉴权
			Metadata(label.Auth, label.Enable).
			// 资源描述
			Metadata(label.Resource, "Secret").
			Metadata(label.Action, label.Create.Value()).
			// 开启鉴权
			Metadata(label.Permission, label.Enable),
	)
```