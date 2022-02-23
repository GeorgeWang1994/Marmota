package msg_opt

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"marmota/alarm/gg"
	"marmota/pkg/common/model"
	"strings"
)

func lpush(queue, message string) {
	rc := gg.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Error("LPUSH redis", queue, "fail:", err, "message:", message)
	}
}

func WriteIMModel(im *model.IM) {
	if im == nil {
		return
	}

	bs, err := json.Marshal(im)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write im to queue, im:%v, queue:%s", im, IM_QUEUE_NAME)
	lpush(IM_QUEUE_NAME, string(bs))
}

func WriteIM(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	im := &model.IM{Tos: strings.Join(tos, ","), Content: content}
	WriteIMModel(im)
}
