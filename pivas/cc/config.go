package cc

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"marmota/pkg/utils/file"
	"sync"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GlobalConfig struct {
	Debug    bool   `json:"debug"`
	Hosts    string `json:"hosts"`
	Database string `json:"database"`
	// 数据库最大连接数
	MaxConns int `json:"maxConns"`
	// 数据库最大空闲数
	MaxIdle   int         `json:"maxIdle"`
	Listen    string      `json:"listen"`
	Trustable []string    `json:"trustable"`
	Http      *HttpConfig `json:"http"`
}

var (
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) error {
	if cfg == "" {
		return errors.New("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return errors.New("config file is not existent")
	}

	configContent, err := file.FileContent(cfg)
	if err != nil {
		return err
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
	return nil
}
