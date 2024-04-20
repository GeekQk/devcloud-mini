package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GeekQk/devcloud-mini/maudit/apps/audit"
	"github.com/infraboard/mcube/v2/ioc/config/kafka"
)

func (i *Consumer) SaveAuditLog() {
	// 消息消费者
	i.l.Debug().Msgf("group id: %s, topic: %s", i.GroupId, i.Topic)
	i.r = kafka.ConsumerGroup(i.GroupId, []string{i.Topic})

	for {
		m, err := i.r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		record := new(audit.AuditLog)
		err = json.Unmarshal(m.Value, record)
		if err != nil {
			i.l.Error().Msgf("unmarshal error, %s", err)
			continue
		}
		err = i.svc.SaveLog(context.Background(), record)
		if err != nil {
			i.l.Error().Msgf("save log error, %s", err)
			continue
		}
	}
}
