package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/set"
	"github.com/toolkits/net/httplib"
	"marmota/alarm/cc"
	"marmota/pkg/common/model"
	"strings"
	"sync"
	"time"
)

type APIGetTeamOutput struct {
	model.Team
	Users       []*model.User `json:"users"`
	TeamCreator string        `json:"creator_name"`
}

type UsersCache struct {
	sync.RWMutex
	M map[string][]*model.User
}

var Users = &UsersCache{M: make(map[string][]*model.User)}

func (u *UsersCache) Get(team string) []*model.User {
	u.RLock()
	defer u.RUnlock()
	val, exists := u.M[team]
	if !exists {
		return nil
	}

	return val
}

func (u *UsersCache) Set(team string, users []*model.User) {
	u.Lock()
	defer u.Unlock()
	u.M[team] = users
}

// CurlUser 通过team获取对应的用户
func CurlUser(team string) []*model.User {
	if team == "" {
		return []*model.User{}
	}

	uri := fmt.Sprintf("%s/api/v1/team/name/%s", cc.Config().Api.PlusApi, team)
	req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  cc.Config().Api.PlusApiToken,
	})
	req.Header("Apitoken", string(token))

	var team_users APIGetTeamOutput
	err := req.ToJson(&team_users)
	if err != nil {
		log.Errorf("curl %s fail: %v", uri, err)
		return nil
	}

	return team_users.Users
}

func UsersOf(team string) []*model.User {
	users := CurlUser(team)

	if users != nil {
		Users.Set(team, users)
	} else {
		users = Users.Get(team)
	}

	return users
}

// GetUsers 获取用户信息，先从缓存中读取，读取不到再通过request去获得
func GetUsers(teams string) map[string]*model.User {
	userMap := make(map[string]*model.User)
	arr := strings.Split(teams, ",")
	for _, team := range arr {
		if team == "" {
			continue
		}

		users := UsersOf(team)
		if users == nil {
			continue
		}

		for _, user := range users {
			userMap[user.Name] = user
		}
	}
	return userMap
}

// ParseTeams return phones, emails, IM
func ParseTeams(teams string) ([]string, []string, []string) {
	if teams == "" {
		return []string{}, []string{}, []string{}
	}

	userMap := GetUsers(teams)
	phoneSet := set.NewStringSet()
	mailSet := set.NewStringSet()
	imSet := set.NewStringSet()
	for _, user := range userMap {
		if user.Phone != "" {
			phoneSet.Add(user.Phone)
		}
		if user.Email != "" {
			mailSet.Add(user.Email)
		}
		if user.IM != "" {
			imSet.Add(user.IM)
		}
	}
	return phoneSet.ToSlice(), mailSet.ToSlice(), imSet.ToSlice()
}
