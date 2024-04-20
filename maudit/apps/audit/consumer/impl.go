package consumer

import (
	"context"

	"github.com/GeekQk/devcloud-mini/maudit/apps/audit"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

func init() {
	ioc.Controller().Registry(&Consumer{
		GroupId: "maudit_consumer",
		Topic:   "audit_log",
	})
}

/*
[audit_consumer]
group_id=xxx
topic=xxx
*/
type Consumer struct {
	// 需要托管Ioc
	ioc.ObjectImpl

	// 一定要定义Json标签
	GroupId string `json:"group_id" toml:"group_id" env:"AUDIT_CONSUMER_GROUP"`
	Topic   string `json:"topic" toml:"topic" env:"AUDIT_CONSUMER_TOPIC"`

	r   *kafka.Reader
	l   *zerolog.Logger
	svc audit.Service
}

// 对象名称
func (i *Consumer) Name() string {
	return "audit_consumer"
}

// 获取consumer
func (i *Consumer) Init() error {
	i.l = log.Sub(i.Name())
	i.svc = ioc.Controller().Get(audit.AppName).(audit.Service)

	// 不能阻塞主流程
	go i.SaveAuditLog()
	return nil
}

func (i *Consumer) Close(ctx context.Context) error {
	if i.r != nil {
		return i.r.Close()
	}
	return nil
}
