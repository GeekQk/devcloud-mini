package audit

import (
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/tools/pretty"
)

type AuditLog struct {
	Who  string `json:"who"`
	When int64  `json:"when"`
	What string `json:"what"`
}

func (s *AuditLog) String() string {
	return fmt.Sprintf("%s %s %s", s.Time(), s.Who, s.What)
}

func (s *AuditLog) Time() time.Time {
	return time.Unix(s.When, 0)
}

func (s *AuditLog) ToJSON() string {
	return pretty.ToJSON(s)
}

type AuditLogSet struct {
	Total int64       `json:"total"`
	Items []*AuditLog `json:"items"`
}
