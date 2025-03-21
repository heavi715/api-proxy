package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/heavi715/api-proxy/public/logger"
)

// Configuration 项目配置
type Configuration struct {
	ServerAddr    string   `json:"server_addr"`
	SourceList    []string `json:"source_list"`
	ServerKeyList []string `json:"server_key_list"`
	ProxyURL      string   `json:"proxy_url"`
	// 调用gpt接口超时时间
	PlatformList map[string]Platform `json:"platform_list"`
}

type Platform struct {
	Name         string   `json:"name"`
	Url          string   `json:"url"`
	HeaderKey    string   `json:"header_key"`
	HeaderValues []string `json:"header_values"`
}

var Config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig(configFile string) *Configuration {

	once.Do(func() {
		// 从文件中读取
		Config = &Configuration{}
		f, err := os.Open(configFile)
		if err != nil {
			logger.Error("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(Config)
		if err != nil {
			logger.Warning("decode config err: %v", err)
			return
		}

	})

	return Config
}

func (config *Configuration) IsSource(source string) bool {
	if config.SourceList == nil {
		return false
	}
	for _, s := range config.SourceList {
		if s == source {
			return true
		}
	}
	return false
}
