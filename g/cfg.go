package g

import (
	"encoding/json"
	"sync"

	"github.com/golang/glog"
	"github.com/toolkits/file"
)

// agent http service config section
type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

// transfer service config section
type TransferConfig struct {
	Enable   bool   `json:"enable"`
	Addr     string `json:"addr"`
	Interval int64  `json:"interval"`
	Timeout  int    `json:"timeout"`
}

// global config file
type GlobalConfig struct {
	Debug      bool            `json:"debug"`
	AttachTags string          `json:"attachtags"`
	Http       *HttpConfig     `json:"http"`
	Transfer   *TransferConfig `json:"transfer"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// parse config file.
func ParseConfig(cfg string) {
	if cfg == "" {
		glog.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		glog.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		glog.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		glog.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	glog.Infoln("g:ParseConfig, ok, ", cfg)
}
