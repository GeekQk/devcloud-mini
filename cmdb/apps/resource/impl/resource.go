package impl

import (
	"context"

	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
)

// 给 secret模块提供的 资源同步方法
func (i *impl) Save(
	ctx context.Context,
	in *resource.Resource) (
	*resource.Resource, error) {

	// 实力数据填充
	in.AutoFill()

	// 保存Resoruce资源, 保存到 resource collection
	_, err := i.col.InsertOne(ctx, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// 对内部系统提供grpc 资源搜索
func (i *impl) Search(
	ctx context.Context,
	in *resource.SearchRequest) (
	*resource.ResourceSet, error) {
	set := resource.NewResourceSet()

	// gorm where 实现条件过滤, mongodb, filter 过滤条件  {name: "xxx"}  where name=xxxx
	cursor, err := i.col.Find(ctx, in.Filter(), in.Options())
	if err != nil {
		return nil, err
	}

	// cursor 是一个迭代器
	// 【a, b, c】 cursor 0, 1, 2
	for cursor.Next(ctx) {
		ins := resource.NewResource()
		// 读取当前Next的值
		if err := cursor.Decode(ins); err != nil {
			return nil, err
		}
		set.Items = append(set.Items, ins)

	}

	return set, nil
}
