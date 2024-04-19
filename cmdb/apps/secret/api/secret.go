package api

import (
	"github.com/GeekQk/devcloud-mini/cmdb/apps/secret"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"
)

// 写一个 go restful 框架的Handler
// 使用的v3版本
// GoRestful (r *restful.Request, w *restful.Response)
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
