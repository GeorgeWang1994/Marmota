package cc

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"marmota/pkg/utils/file"
)

var (
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

type ClusterNode struct {
	Addrs []string `json:"addrs"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int    `json:"timeout"`
}

type TsdbConfig struct {
	Enabled     bool   `json:"enabled"`
	Batch       int    `json:"batch"`
	ConnTimeout int    `json:"connTimeout"`
	CallTimeout int    `json:"callTimeout"`
	MaxConns    int    `json:"maxConns"`
	MaxIdle     int    `json:"maxIdle"`
	MaxRetry    int    `json:"retry"`
	Address     string `json:"address"`
}

type JudgeConfig struct {
	Enabled     bool                    `json:"enabled"`
	Batch       int                     `json:"batch"`
	ConnTimeout int                     `json:"connTimeout"`
	CallTimeout int                     `json:"callTimeout"`
	MaxConns    int                     `json:"maxConns"`
	MaxIdle     int                     `json:"maxIdle"`
	Replicas    int                     `json:"replicas"`
	Cluster     map[string]string       `json:"cluster"`
	ClusterList map[string]*ClusterNode `json:"clusterList"`
}

type GraphConfig struct {
	Enabled     bool                    `json:"enabled"`
	Batch       int                     `json:"batch"`
	ConnTimeout int                     `json:"connTimeout"`
	CallTimeout int                     `json:"callTimeout"`
	MaxConns    int                     `json:"maxConns"`
	MaxIdle     int                     `json:"maxIdle"`
	Replicas    int                     `json:"replicas"`
	Cluster     map[string]string       `json:"cluster"`
	ClusterList map[string]*ClusterNode `json:"clusterList"`
}

type TransferConfig struct {
	Enabled     bool              `json:"enabled"`
	Batch       int               `json:"batch"`
	ConnTimeout int               `json:"connTimeout"`
	CallTimeout int               `json:"callTimeout"`
	MaxConns    int               `json:"maxConns"`
	MaxIdle     int               `json:"maxIdle"`
	MaxRetry    int               `json:"retry"`
	Cluster     map[string]string `json:"cluster"`
}

type InfluxdbConfig struct {
	Enabled   bool   `json:"enabled"`
	Batch     int    `json:"batch"`
	MaxRetry  int    `json:"retry"`
	MaxConns  int    `json:"maxConns"`
	Timeout   int    `json:"timeout"`
	Address   string `json:"address"`
	Database  string `json:"db"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Precision string `json:"precision"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	MinStep  int             `json:"minStep"` //最小周期,单位sec
	Http     *HttpConfig     `json:"http"`
	Rpc      *RpcConfig      `json:"connpool"`
	Socket   *SocketConfig   `json:"socket"`
	Judge    *JudgeConfig    `json:"judge"`
	Graph    *GraphConfig    `json:"graph"`
	Tsdb     *TsdbConfig     `json:"tsdb"`
	Transfer *TransferConfig `json:"transfer"`
	Influxdb *InfluxdbConfig `json:"influxdb"`
}

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
