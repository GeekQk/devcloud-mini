# 资源同步模块


## secret管理

```go
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

	return ins, nil
}
```

## 加密存储

[](./docs/desense.drawio)

```go
// 加密
// 非对称加密: 加密密钥 (密钥对)
// 对称加密: 数据 (key: "xxx")
// 消息摘要(Hash)
// key 通过配置保存
// passwd   cbc.v1.abc.<cipherText>
func (s *Secret) Encrypt() error {
	cipherText, err := cbc.EncryptToString(s.Spec.Value, []byte(application.Get().EncryptKey))
	if err != nil {
		return err
	}
	s.Spec.Value = cipherText
	return nil
}

func (s *Secret) Decrypt() error {
	planText, err := cbc.DecryptFromString(s.Spec.Value, []byte(application.Get().EncryptKey))
	if err != nil {
		return err
	}
	s.Spec.Value = planText
	return nil
}

func (s *Secret) Desense() {
	// api secret
	// api****ret
	s.Spec.Value = s.Spec.Value[:3] + "****"
}
```


## 整合Provider模块 进行资源同步

[](../../provider/sync.drawio)


```go
func TestSearchResource(t *testing.T) {
	req := resource.NewSearchRequest()
	req.Keywords = "宝塔Linux面板"
	res, err := impl.Search(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
```

核心同步逻辑
```go
	// 3. 使用 CVMProvider 进行资源同步
	// 同步 进行资源同步，还是异步进行资源同步 (task)
	// ctx 请求的context, 任务还需要继续进行
	gctx, _ := context.WithTimeout(context.Background(), 1*time.Hour)
	for _, rp := range resourceProviders {
		// 并发异步同步
		go rp.Sync(gctx,
			func(ctx context.Context, r *resource.Resource) {
				_, err := i.resource.Save(ctx, r)
				if err != nil {
					// 同步失败
					log.L().Error().Msgf("save resource error, %s", err)
				} else {
					// 同步成功
					rh(&secret.SyncResponse{
						Id:   r.Meta.Id,
						Name: r.Spec.Name,
					})
				}
			})
	}
```


## upsert

新增或者更新

```go
// 给 secret模块提供的 资源同步方法
func (i *impl) Save(
	ctx context.Context,
	in *resource.Resource) (
	*resource.Resource, error) {

	// 实力数据填充
	in.AutoFill()

	// 保存Resoruce资源, 保存到 resource collection

	// 如果没有则创建 如果有责更新
	op := options.Update().SetUpsert(true)

	_, err := i.col.UpdateOne(ctx, bson.M{"_id": in.Meta.Id}, bson.M{"$set": in}, op)
	if err != nil {
		return nil, err
	}

	return in, nil
}
```

## api

```go
func (h *handler) SyncResource(r *restful.Request, w *restful.Response) {
	req := &secret.SyncResourceRequest{}
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	req.SecretId = r.PathParameter("id")
	err := h.secret.SyncResource(r.Request.Context(), req, func(sr *secret.SyncResponse) {
		h.log.Debug().Msgf("%s[%s], %s %s", sr.Name, sr.Id, sr.Status(), sr.Error)
	})
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, "ok")
}
```

```sh
curl --location 'http://127.0.0.1:7010/cmdb/api/v1/secret/cod36eh97i69ji730se0/sync' \
--header 'Content-Type: application/json' \
--header 'Cookie: mcenter.access_token=aLTX3EdUFlqBbHjeU9akG169' \
--data '{
    
}'
```

## 关于定时认证


[](./docs/sync.drawio)

+ 自己集成cron功能: https://github.com/robfig/cron
```go
	cron.AddFunc("0 0 0 1 1 ?", func() {})
	cron.AddFunc("0 0 0 31 12 ?", func() {})
	cron.AddFunc("* * * * * ?", func() { wg.Done() })
```

程序会有状态, 每个实例都会执行同步

+ 使用外部Job管理平台: (crontab, k8s cronjob, cronjob管理平台)



## 优化


### mongodb 数据存储优化

```proto
message Secret {
    // secret id
    // @gotags: json:"id" bson:"_id"
    string id = 1;
    // 创建secret的请求 
    // @gotags: json:"spec" bson:",inline"
    CreateSecretRequest spec = 4;
}
```


### api 接口返回优化


