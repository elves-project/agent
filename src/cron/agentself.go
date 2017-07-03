package cron

import (
	"../funcs"
	"../g"
	"../thrift/scheduler"
	"encoding/json"
	"fmt"
	"github.com/gy-games-libs/cron"
	"github.com/gy-games-libs/file"
	"github.com/gy-games-libs/fsnotify"
	"github.com/gy-games-libs/seelog"
	"strings"
	"sync"
)

type CLStrct struct {
	Flag     string            `json:"flag"`
	Comment  string            `json:"comment"`
	Id       string            `json:"id"`
	Mode     string            `json:"mode"`
	App      string            `json:"app"`
	Func     string            `json:"func"`
	Param    map[string]string `json:"param"`
	Timeout  int32             `json:"timeout"`
	Proxy    string            `json:"proxy"`
	Rule     string            `json:"rule"`
	LastExec string
}

var (
	PATH     string
	CronList map[string]CLStrct
	CRON     = cron.New()
	RUNJOB   = map[string]string{}
	lock     = new(sync.Mutex)
	clock    = new(sync.Mutex)
)

func WatchCron() {
	seelog.Info("watch agent cron")
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("[funcs:WatchCron] ", err)
			go g.SaveErrorStat("[func:WatchCron] " + fmt.Sprintf("%s", err))
		}
	}()
	PATH = g.Root + "/conf/cron.json"
	seelog.Debug("[func:WatchCron] ", PATH)
	Watch, err := fsnotify.NewWatcher()
	if err != nil {
		seelog.Error("[func:WatchCron] Init monitor error: ", err.Error())
		return
	}
	if err := Watch.Add(PATH); err != nil {
		seelog.Error("[func:WatchCron] Add monitor path error: ", PATH)
		return
	}
	Updatecron()
	CRON.Start()
	for {
		select {
		case event := <-Watch.Events:
			seelog.Debug("[func:WatchCron] ", event.Op.String())
			if event.Op.String() == "WRITE" || event.Op.String() == "REMOVE" {
				lock.Lock()
				seelog.Info("Cron File Changed Update Now Start..")
				Updatecron()
				seelog.Info("Cron Update Finish..")
				lock.Unlock()
				if event.Op.String() == "REMOVE" {
					if err := Watch.Add(PATH); err != nil {
						seelog.Error("[func:WatchCron] ReAdd monitor path error: ", PATH)
						return
					}
				}
			}
		case err := <-Watch.Errors:
			seelog.Error("[func:WatchCron] ", err)
			go g.SaveErrorStat("[func:WatchCron] " + err.Error())
			//return
		}
	}
}

func Updatecron() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("[funcs:UpdateCron] ", err)
			go g.SaveErrorStat("[func:UpdateCron] " + fmt.Sprintf("%s", err))
		}
	}()
	defer clock.Unlock()
	clock.Lock()
	configContent, err := file.ToTrimString(PATH)
	if err != nil {
		seelog.Error("[func:WatchCron] read cron config file: ", PATH, "fail:", err)
		go g.SaveErrorStat("[func:WatchCron] read cron config file: " + PATH + " fail:" + err.Error())
	}
	err = json.Unmarshal([]byte(configContent), &CronList)
	if err != nil {
		seelog.Error("[func:WatchCron] parse cron config file:", PATH, "fail:", err)
		go g.SaveErrorStat("[func:WatchCron] parse cron config file:" + PATH + " fail:" + err.Error())
	}
	for k, cronjob := range CronList {
		seelog.Debug("[func:WatchCron] k:", k, " v:", cronjob)
		ins := new(scheduler.Instruct)
		Param, err := json.Marshal(cronjob.Param)
		ins.Param = funcs.BytesString(Param)
		if err == nil {
			ins.App = cronjob.App
			ins.Func = cronjob.Func
			ins.Type = "cron"
			ins.ID = k
			ins.Mode = strings.ToUpper(cronjob.Mode)
			ins.IP = g.Config().Ip
			ins.Proxy = cronjob.Proxy
			ins.Timeout = cronjob.Timeout
			md5check, _ := json.Marshal(ins.App + ins.Func + ins.ID + ins.Mode + ins.Proxy + ins.Param + string(ins.Timeout) + cronjob.Rule + string(cronjob.Flag) + cronjob.Comment)
			m := funcs.GetMD5(string(md5check))
			seelog.Debug("[func:WatchCron] ID:" + ins.ID + " MD5:" + m)
			if RUNJOB[ins.ID] != m || RUNJOB[ins.ID] == "" {
				if RUNJOB[ins.ID] != "" {
					CRON.RemoveJob(ins.ID)
					seelog.Debug("[func:WatchCron] Remove Cron..", cronjob)
					g.DelCronSata(ins.ID)
				}
				if cronjob.Flag == "true" {
					CRON.AddFunc(cronjob.Rule, func() {
						funcs.ExecAndSend(*ins, false)
						g.UpdateCronSata(ins.ID)
					}, ins.ID)
					RUNJOB[ins.ID] = m
					g.SaveCronSata(ins.ID, ins.App, ins.Func, ins.Mode, cronjob.Rule, cronjob.Comment)
					seelog.Debug("[func:WatchCron] Add Cron..", cronjob)
				}
			}
		} else {
			seelog.Error("[func:WatchCron] ParamError ", err)
			go g.SaveErrorStat("[func:WatchCron] ParamError " + err.Error())
		}
		for rk, _ := range RUNJOB {
			if CronList[rk].Flag == "false" || CronList[rk].Flag == "" {
				CRON.RemoveJob(rk)
				g.DelCronSata(rk)
				RUNJOB[ins.ID] = ""
				seelog.Debug("[func:WatchCron] Remove Cron..", rk)
			}
		}
		seelog.Debug("[func:WatchCron] CRON Jobs ", RUNJOB)
	}
	seelog.Debug("[func:WatchCron] ", CronList)
}
