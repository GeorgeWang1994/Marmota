package consume

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"marmota/alarm/api"
	"marmota/alarm/cc"
	"marmota/alarm/cron/event/msg_opt"
	"marmota/alarm/gg"
	"marmota/pkg/common/model"
)

func consume(event *model.Event, isHigh bool) {
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	action := api.GetAction(actionId)
	if action == nil {
		return
	}

	if isHigh {
		consumeHighEvents(event, action)
	} else {
		consumeLowEvents(event, action)
	}
}

// 高优先级的不做报警合并
func consumeHighEvents(event *model.Event, action *model.Action) {
	_, _, ims := api.ParseTeams(action.UIC)

	imContent := GenerateIMContent(event)

	msg_opt.WriteIM(ims, imContent)
}

// 低优先级的做报警合并
func consumeLowEvents(event *model.Event, action *model.Action) {
	ParseUserIm(event, action)
}

// ParseUserIm 即使消息
func ParseUserIm(event *model.Event, action *model.Action) {
	userMap := api.GetUsers(action.UIC)

	content := GenerateIMContent(event)
	metric := event.Metric()
	status := event.Status
	priority := event.Priority()

	queue := cc.Config().Redis.UserIMQueue

	rc := gg.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := ImDto{
			Priority: priority,
			Metric:   metric,
			Content:  content,
			IM:       user.IM,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Error("json marshal ImDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Error("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}
