package audit

import "context"

const (
	AppName = "audit_log"
)

type Service interface {
	// 查询接口
	QueryLog(context.Context, *QueryLogRequest) (*AuditLogSet, error)
	// 录入
	SaveLog(context.Context, *AuditLog) error
}

type QueryLogRequest struct {
}
