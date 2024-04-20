package impl

import (
	"context"

	"github.com/GeekQk/devcloud-mini/maudit/apps/audit"
)

// 查询接口
func (i *impl) QueryLog(
	ctx context.Context,
	in *audit.QueryLogRequest,
) (*audit.AuditLogSet, error) {
	return nil, nil
}

// 录入
func (i *impl) SaveLog(
	ctx context.Context,
	in *audit.AuditLog) error {
	i.log.Debug().Msgf("save log: %s", in)
	return nil
}
