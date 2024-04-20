# 审计需求

1. 用户操作审计, 避免用户赖账
2. 给安全团队做安全分析


## 需求

核心需求: WWW原则
1. Who:    token --->  username, 做在认证中间件后面
1. When:  接口调用事件
1. What:  接口名称, 接口的一些核心参数(资源ID)


## 架构设计

产品架构:

+ 审计日志(审计管理员): 

数据架构:

+ 审计数据的录入: 需要接入的服务
+ 数据查询: UI

[](./data_arch.drawio)

服务架构:
+ 多个服务之间的管理: maudit

## 开发

+ api: gorestful 开发
+ 发送日志: kafka的 writer, 发生日志
+ 接收日志: kafka的 reader