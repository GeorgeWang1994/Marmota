package cc

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"marmota/pkg/utils/file"
	"sync"
)

var (
	lock sync.RWMutex
	cc   *GlobalConfig
)

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type HbsConfig struct {
	Servers  []string `json:"servers"`
	Timeout  int64    `json:"timeout"`
	Interval int64    `json:"interval"`
}

type RedisConfig struct {
	Dsn          string `json:"dsn"`
	MaxIdle      int    `json:"maxIdle"`
	ConnTimeout  int    `json:"connTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

type AlarmConfig struct {
	Enabled      bool         `json:"enabled"`
	MinInterval  int64        `json:"minInterval"`  // 告警最小间隔事件
	QueuePattern string       `json:"queuePattern"` // 告警的队列，用来告诉发往哪个redis key
	Redis        *RedisConfig `json:"redis"`
}

type GlobalConfig struct {
	Debug     bool         `json:"debug"`
	DebugHost string       `json:"debugHost"`
	Remain    int          `json:"remain"`
	Rpc       *RpcConfig   `json:"rpc"`
	Hbs       *HbsConfig   `json:"hbs"`
	Alarm     *AlarmConfig `json:"alarm"`
}

func ParseConfig(cfg string) error {
	if cfg == "" {
		return errors.New("use -c to specify pivas file")
	}

	if !file.IsExist(cfg) {
		return errors.New("is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	configContent, err := file.FileContent(cfg)
	if err != nil {
		return err
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return err
	}

	lock.Lock()
	defer lock.Unlock()

	cc = &c

	log.Println("read cc file:", cfg, "successfully")
	return nil
}

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return cc
}
