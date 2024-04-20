package middleware

import (
	"time"

	"github.com/GeekQk/devcloud-mini/maudit/apps/audit"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	ioc_kafka "github.com/infraboard/mcube/v2/ioc/config/kafka"
	"github.com/infraboard/mcube/v2/ioc/config/log"
)

func init() {
	ioc.Config().Registry(&AuditLogProducer{
		Topic: "audit_log",
	})
}

/*
[audit_log_producer]
topic = "xx"
*/
// 需要ioc里面获取 producer对象
type AuditLogProducer struct {
	ioc.ObjectImpl
	// 一定要定义Json标签
	Topic string `json:"topic" toml:"topic" env:"AUDIT_CONSUMER_TOPIC"`

	l *zerolog.Logger
}

func (p *AuditLogProducer) Name() string {
	return "audit_log_producer"
}

func (p *AuditLogProducer) Init() error {
	p.l = log.Sub(p.Name())

	// 怎么中间件加入到router
	gorestful.RootRouter().Filter(p.SendFunc)
	return nil
}

func (p *AuditLogProducer) SendFunc(
	req *restful.Request,
	resp *restful.Response,
	next *restful.FilterChain) {
	// 只有开启了认证中间件的 才能获取用户是谁
	tk := token.GetTokenFromRequest(req)
	if tk != nil {
		// 1. record
		record := audit.AuditLog{
			Who:  tk.Username,
			When: time.Now().Unix(),
			What: req.SelectedRoute().Operation(),
		}

		// 生产日志
		producer := ioc_kafka.Producer(p.Topic)
		err := producer.WriteMessages(req.Request.Context(), kafka.Message{
			Value: []byte(record.ToJSON()),
		})
		if err != nil {
			p.l.Debug().Msgf("send audit log error, %s", err)
		}
	}

	next.ProcessFilter(req, resp)
}
