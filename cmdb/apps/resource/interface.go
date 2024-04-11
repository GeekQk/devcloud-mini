package resource

import (
	context "context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AppName = "resource"
)

// 资源录入是 cmdb的secret提供的功能, 并不对外通过rpc开发资源录入的功能
type Service interface {
	// 资源保持接口
	Save(context.Context, *Resource) (*Resource, error)
	RPCServer
}

func NewSearchRequest() *SearchRequest {
	return &SearchRequest{
		PageSize:   10,
		PageNumber: 1,
		Tags:       []*TagSelector{},
	}
}

// Mongo的 查询过滤条件
// https://www.mongodb.com/docs/manual/reference/command/find/
func (r *SearchRequest) Filter() bson.M {
	filter := bson.M{}

	// 精确调节匹配
	if r.Type != nil {
		// resource {spec:{type: "host"}}
		filter["spec.type"] = r.Type.String()
	}

	// 模糊匹配, 使用正则做匹配
	// LIKE name = %?%
	// $regex 条件:  后面正则的值
	// $gte >=
	// $gt >
	if r.Keywords != "" {
		filter["spec.name"] = bson.M{"$regex": r.Keywords, "$options": "im"}
	}

	return filter
}

// monogo db 的查询分页 Limit offset,limite
func (r *SearchRequest) Options() *options.FindOptions {
	opt := options.Find()

	// 忽略前面多少个
	opt.SetSkip(r.Offset())
	opt.SetLimit(int64(r.PageSize))
	return nil
}

func (r *SearchRequest) Offset() int64 {
	return int64((r.PageSize - 1) * r.PageSize)
}
