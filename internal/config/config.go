package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configsPath = "configs"

var (
	configFilePath = fmt.Sprintf("%s/%s", configsPath, "config.json")
	Config         config
)

type config struct {
	Host            string   `json:"host"`
	Port            uint16   `json:"port"`
	UpdateAt        uint32   `json:"updateAt"`
	Nodes           []string `json:"nodes"`
	SubscribeLink   []string `json:"subscribeLink"`
	SubscribeBase64 []string `json:"subscribeBase64"`
}

func (c *config) ReadConfig() error {
	var err error
	var file []byte
	if file, err = os.ReadFile(configFilePath); err == nil {
		if err = json.Unmarshal(file, c); err != nil {
			log.Printf("已成功读取配置文件")
		}
	}
	return err
}

func (c *config) InitConfig() error {
	var err error
	c.Host = "127.0.0.1"
	c.Port = 38080
	c.UpdateAt = 14400
	c.Nodes = []string{}
	c.SubscribeLink = []string{}
	c.SubscribeBase64 = []string{}

	var file []byte
	if file, err = json.MarshalIndent(c, "", "    "); err == nil {
		if err = os.WriteFile(configFilePath, file, os.FileMode(0660)); err == nil {
			log.Printf("已成功写入初始化配置文件")
		}
	}
	return err
}

func init() {
	var err error
	os.MkdirAll(configsPath, os.FileMode(0660))
	if err = Config.ReadConfig(); err != nil {
		if err = Config.InitConfig(); err != nil {
			if err = Config.ReadConfig(); err != nil {
				log.Fatalf("配置文件出现无法处理的异常, 请重试: %v", err)
			}
		}
	}
}
