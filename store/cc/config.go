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

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RRDConfig struct {
	Storage string `json:"storage"`
}

type DBConfig struct {
	Dsn     string `json:"dsn"`
	MaxIdle int    `json:"maxIdle"`
}

type GlobalConfig struct {
	Pid            string      `json:"pid"`
	Debug          bool        `json:"debug"`
	Http           *HttpConfig `json:"http"`
	Rpc            *RpcConfig  `json:"rpc"`
	RRD            *RRDConfig  `json:"rrd"`
	DB             *DBConfig   `json:"db"`
	CallTimeout    int32       `json:"callTimeout"`
	IOWorkerNum    int         `json:"ioWorkerNum"`
	FirstBytesSize int
	Migrate        struct {
		Concurrency int               `json:"concurrency"` //number of multiple worker per node
		Enabled     bool              `json:"enabled"`
		Replicas    int               `json:"replicas"`
		Cluster     map[string]string `json:"cluster"`
	} `json:"migrate"`
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
