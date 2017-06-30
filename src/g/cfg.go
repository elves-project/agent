package g

import (
	"github.com/gy-games-libs/seelog"
	"github.com/gy-games-libs/file"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"log"
	"io/ioutil"
	"sync"
	"encoding/json"
)

type AgentCronConfig struct{
	Comment 		string			`json:"comment"`
	Enabled 		bool			`json:"enabled"`
}

type HeartbeatConfig struct {
	Comment 		string			`json:"comment"`
	Enabled 		bool			`json:"enabled"`
	Addr			string			`json:"addr"`
	Port 	 		int			`json:"port"`
	Interval 		int			`json:"interval"`
	Timeout  		int			`json:"timeout"`
}

type SchedulerConfig struct{
	Addr 	string			`json:"addr"`
	Port 	int			`json:"port"`
	Timeout	int			`json:"timeout"`
}

type DevmodeConfig struct {
	Comment 		string			`json:"comment"`
	Enabled 		bool			`json:"enabled"`
	Authid			string			`json:"authid"`
	Authkey 	 	string			`json:"authkey"`
}

type GlobalConfig struct{
	Asset			string			`json:"asset"`
	Ip 			string			`json:"ip"`
	Port            	int			`json:"port"`
	HeartBeat		*HeartbeatConfig	`json:"heartbeat"`
	Scheduler		*SchedulerConfig	`json:"scheduler"`
	Agentcron		*AgentCronConfig	`json:"agentcron"`
	Http			*HttpConfig		`json:"http"`
	AppsDownloadAddr	string			`json:"appsdownloadaddr"`
	Apps			map[string]string	`json:"apps"`
	Devmode			*DevmodeConfig		`json:"devmode"`
}

type HttpConfig struct {
	Comment 		string			`json:"comment"`
	Enabled 		bool			`json:"enabled"`
	Listen            	string			`json:"listen"`
}

var (
	Root 		string
	ConfigFile 	string
	config     	*GlobalConfig
	lock       	= new(sync.RWMutex)
	savecfglock 	= new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func GetRoot() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dirctory := strings.Replace(dir, "\\", "/", -1)
	runes := []rune(dirctory)
	l := 0 + strings.LastIndex(dirctory, "/")
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[0:l])
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -r to specify elves-agent root directory")
	}

	if !file.IsExist(cfg+"/conf/"+"cfg.json") {
		log.Fatalln("[Fault]config file:", cfg+"/conf/"+"cfg.json", "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	Root = cfg

	configContent, err := file.ToTrimString(cfg+"/conf/"+"cfg.json")
	if err != nil {
		log.Fatalln("[Fault]read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("[Fault]parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	if config.Ip ==""{
		log.Fatalln(cfg,"[Fault]get local ip addr fail check you config!")
	}

	if config.Asset == "" {
		config.Asset = config.Ip
	}

	logger, err := seelog.LoggerFromConfigAsFile(cfg+"/conf/"+"seelog.xml")

	seelog.ReplaceLogger(logger)

}

func SaveConfig()  {
	savecfglock.RLock()
	defer savecfglock.RUnlock()
	seelog.Debug("Save Config..",Root+"/conf/"+"cfg.json")
	r,_ := json.Marshal(Config())
	var out bytes.Buffer
	err := json.Indent(&out, r, "", "\t")

	if err != nil {
		seelog.Error(err)
	}
	if err == nil {
		//seelog.Debug(out.String())
		ioutil.WriteFile(Root+"/conf/"+"cfg.json", []byte(out.String()),0644)
		seelog.Debug("Config Save Success!")
	}else{
		seelog.Error("save config fail ",err)
	}
}