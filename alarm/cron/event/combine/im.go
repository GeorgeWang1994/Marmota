package combine

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"marmota/alarm/cc"
	"marmota/alarm/cron/event/consume"
	"marmota/alarm/cron/event/msg_opt"
	"marmota/alarm/gg"
	"strings"
	"time"
)

func CombineIM() {
	for {
		// 每分钟读取处理一次
		time.Sleep(time.Minute)
		combineIM()
	}
}

func popAllImDto() []*consume.ImDto {
	ret := []*consume.ImDto{}
	queue := cc.Config().Redis.UserIMQueue

	rc := gg.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get ImDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var imDto consume.ImDto
		err = json.Unmarshal([]byte(reply), &imDto)
		if err != nil {
			log.Errorf("json unmarshal imDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &imDto)
	}

	return ret
}

func combineIM() {
	dtos := popAllImDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*consume.ImDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].IM, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*consume.ImDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			msg_opt.WriteIM([]string{arr[0].IM}, arr[0].Content)
			continue
		}

		// 把多个im内容写入数据库，只给用户提供一个链接
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, ",,")
		msg_opt.WriteIM([]string{arr[0].IM}, content)
	}
}
