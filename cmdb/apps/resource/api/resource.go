package api

import (
	"github.com/GeekQk/devcloud-mini/cmdb/apps/resource"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
)

func (h *handler) SearchResource(r *restful.Request, w *restful.Response) {
	req := resource.NewSearchRequest()

	set, err := h.resource.Search(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 构建一个临时对象
	warpSet := &ResourceSet{
		Items: []*Resource{},
		Total: set.Total,
	}

	for _, item := range set.Items {
		warpSet.Items = append(warpSet.Items, &Resource{
			Meta:   item.Meta,
			Spec:   item.Spec,
			Status: item.Status,
		})
	}

	response.Success(w, warpSet)
}

type ResourceSet struct {
	Items []*Resource `json:"items"`
	Total int64       `json:"total"`
}

type Resource struct {
	// 资源元数据信息
	*resource.Meta
	// 资源规格信息
	*resource.Spec
	// 资源状态
	Status *resource.Status `json:"status"`
}
