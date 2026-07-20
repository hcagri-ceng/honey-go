package models

import (
	"time"
)

type Event struct {
	ID         int64
	Timestamp  time.Time
	SourceIP   string
	SourcePort int
	TargetPort int
	Protocol   string
	Payload    []byte
}
