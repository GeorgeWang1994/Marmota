package model

import (
	"fmt"
	"marmota/pkg/utils/date"
	"marmota/pkg/utils/format"
)

// Event 机器监控和实例监控都会产生Event，共用这么一个struct
type Event struct {
	Id          string            `json:"id"`
	Strategy    *Strategy         `json:"strategy"`
	Expression  *Expression       `json:"expression"`
	Status      string            `json:"status"`   // OK or PROBLEM
	Endpoint    string            `json:"endpoint"` //
	LeftValue   float64           `json:"leftValue"`
	CurrentStep int               `json:"currentStep"` // 当前告警次数
	EventTime   int64             `json:"eventTime"`   // 产生事件的事件点
	PushedTags  map[string]string `json:"pushedTags"`
}

func (e *Event) FormattedTime() string {
	return date.UnixTsFormat(e.EventTime)
}

func (e *Event) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Status:%s, Strategy:%v, Expression:%v, LeftValue:%s, CurrentStep:%d, PushedTags:%v, TS:%s>",
		e.Endpoint,
		e.Status,
		e.Strategy,
		e.Expression,
		format.ReadableFloat(e.LeftValue),
		e.CurrentStep,
		e.PushedTags,
		e.FormattedTime(),
	)
}

func (e *Event) Priority() int {
	if e.Strategy != nil {
		return e.Strategy.Priority
	}
	return e.Expression.Priority
}
