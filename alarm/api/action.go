package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net/httplib"
	"marmota/alarm/cc"
	"marmota/pkg/common/model"
	"sync"
	"time"
)

type ActionCache struct {
	sync.RWMutex
	M map[int]*model.Action
}

var Actions = &ActionCache{M: make(map[int]*model.Action)}

func (this *ActionCache) Get(id int) *model.Action {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[id]
	if !exists {
		return nil
	}

	return val
}

func (this *ActionCache) Set(id int, action *model.Action) {
	this.Lock()
	defer this.Unlock()
	this.M[id] = action
}

func GetAction(id int) *model.Action {
	action := CurlAction(id)

	if action != nil {
		Actions.Set(id, action)
	} else {
		action = Actions.Get(id)
	}

	return action
}

func CurlAction(id int) *model.Action {
	if id <= 0 {
		return nil
	}

	uri := fmt.Sprintf("%s/api/v1/action/%d", cc.Config().Api.PlusApi, id)
	req := httplib.Get(uri).SetTimeout(5*time.Second, 30*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  cc.Config().Api.PlusApiToken,
	})
	req.Header("Apitoken", string(token))

	var act model.Action
	err := req.ToJson(&act)
	if err != nil {
		log.Errorf("curl %s fail: %v", uri, err)
		return nil
	}

	return &act
}
