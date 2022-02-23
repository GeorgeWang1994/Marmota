package consume

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net/httplib"
	"marmota/alarm/cc"
	"marmota/alarm/cron/event/msg_opt"
	"marmota/pkg/common/model"
	"time"
)

var (
	IMWorkerChan chan int
)

func ConsumeIM() {
	for {
		L := msg_opt.PopAllIM()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendIMList(L)
	}
}

func SendIMList(L []*model.IM) {
	for _, im := range L {
		IMWorkerChan <- 1
		go SendIM(im)
	}
}

func SendIM(im *model.IM) {
	defer func() {
		<-IMWorkerChan
	}()

	url := cc.Config().Api.IM
	r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	r.Param("tos", im.Tos)
	r.Param("content", im.Content)
	resp, err := r.String()
	if err != nil {
		log.Errorf("send im fail, tos:%s, content:%s, error:%v", im.Tos, im.Content, err)
	}

	log.Debugf("send im:%v, resp:%v, url:%s", im, resp, url)
}
