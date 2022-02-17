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

type WorkerConfig struct {
	IM   int `json:"im"`
	Sms  int `json:"sms"`
	Mail int `json:"mail"`
}

type ApiConfig struct {
	Sms          string `json:"sms"`
	Mail         string `json:"mail"`
	Dashboard    string `json:"dashboard"`
	PlusApi      string `json:"plus_api"`
	PlusApiToken string `json:"plus_api_token"`
	IM           string `json:"im"`
}

type FalconPortalConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RedisConfig struct {
	Dsn           string   `json:"dsn"`
	MaxIdle       int      `json:"maxIdle"`
	ConnTimeout   int      `json:"connTimeout"`
	ReadTimeout   int      `json:"readTimeout"`
	WriteTimeout  int      `json:"writeTimeout"`
	HighQueues    []string `json:"highQueues"`
	LowQueues     []string `json:"lowQueues"`
	UserIMQueue   string   `json:"userIMQueue"`
	UserSmsQueue  string   `json:"userSmsQueue"`
	UserMailQueue string   `json:"userMailQueue"`
}

type GlobalConfig struct {
	LogLevel     string              `json:"log_level"`
	FalconPortal *FalconPortalConfig `json:"falcon_portal"`
	Http         *HttpConfig         `json:"http"`
	Redis        *RedisConfig        `json:"redis"`
	Api          *ApiConfig          `json:"api"`
	Worker       *WorkerConfig       `json:"worker"`
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
