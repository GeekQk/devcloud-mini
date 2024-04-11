# cmdb 产品需求

正在传统的cmdb 是一个模型库

没有cmdb之前, 有很多信息是保存在excel:
+ 主机信息
+ IDC 信息
+ 应用信息(地址, 负责人)
+ 数据库:, 链接地址, 用户, 密码
+ 证书: 
+ ....

cmdb相当于一个数据库, 创建模型， 类比为一个class, Struct
```go
// Host 模型
type Host struct {
    IP: xxx
    Name: xxxx
    IDC: xxx
}

// 添加Host 实例
ip: 10.10.1.20 name: 开发a使用, idc: 北京
ip: 10.10.1.21 name: 开发a使用, idc: 北京
ip: 10.10.1.22 name: 开发a使用, idc: 北京
```

理解为:  
+ 数据库里面创建: Host,  把excel 主机Host的信息都放到这个表里面
+ 数据库里面创建: Idc,  把excel 主机IDC的信息都放到这个表里面
+ ...

CMDB:  Config Management 数据库


对外提供Api接口 查看表里面的数据
+ Host.List()
+ App.List()
+ Cert.List() 

有了接口, 重要化的数据库, 客供外部系统集成

CMDB系统核心诉求:
+ 数据自动录入(数据维护的效率)
+ 通过过API 为其他系统提供数据集成能力

[](./arch.drawio)

## CMDB的方案

[](./design.drawio)


