package cron

import (
	"../funcs"
	"../g"
	"../thrift/hearbeat"
	"github.com/gy-games-libs/seelog"
	"github.com/gy-games-libs/go-thrift"
	"os"
	"strconv"
	"net"
	"encoding/json"
	"time"
)

type HeartBeatMessage struct {
	Data     map[string]string	`json:"data"`
}

func HearBeatCron(sec int64){
	t:= time.NewTicker(time.Second*time.Duration((sec))).C
	for{
		<-t
		go sendToHeartBeat()
		go g.UpdateHbTime()
	}
}

func sendToHeartBeat() {
	seelog.Info("Send Info To HeartBeat Start...")
	cfg := g.Config()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	seelog.Debug("[func:SendToHeartBeat] HeartBeat Addr:"+cfg.HeartBeat.Addr+":"+strconv.Itoa(cfg.HeartBeat.Port))
	transport, err := thrift.NewTSocket(net.JoinHostPort(cfg.HeartBeat.Addr, strconv.Itoa(cfg.HeartBeat.Port)))
	defer transport.Close()
	transport.SetTimeout(time.Duration(cfg.HeartBeat.Timeout)*time.Second)
	if err != nil {
		seelog.Error("[func:SendToHeartBeat] ",os.Stderr, "Error resolving address:", err)
		go g.SaveErrorStat("[func:SendToHeartBeat] "+ "Error resolving address:"+ err.Error())
	}else{
		client := hearbeat.NewHeartbeatServiceClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			seelog.Error("[func:SendToHeartBeat] ",os.Stderr, "Error opening socket to "+cfg.HeartBeat.Addr+":"+string(cfg.HeartBeat.Port), " ", err)
			go g.SaveErrorStat("[func:SendToHeartBeat] "+ "Error opening socket to "+cfg.HeartBeat.Addr+":"+string(cfg.HeartBeat.Port)+ " "+ err.Error())
		}else{
			ai := &hearbeat.AgentInfo{}
			ai.IP = cfg.Ip
			ai.ID = cfg.Asset
			ai.Version = g.VERSION
			Apps,_ := json.Marshal(cfg.Apps)
			ai.Apps = funcs.BytesString(Apps)
			jsrets,hberr := client.HeartbeatPackage(ai)
			if hberr == nil {
				var rdat HeartBeatMessage
				seelog.Debug("[func:SendToHeartBeat] HeartBeat Return App " , jsrets)
				json.Unmarshal([]byte(jsrets), &rdat)
				if err := json.Unmarshal([]byte(jsrets), &rdat); err == nil {
					for b,v :=  range rdat.Data{
						if _, ok := cfg.Apps[b]; ok {
							if v != cfg.Apps[b] {
								seelog.Debug("[func:SendToHeartBeat] App '" + b + "' Need Update(Local Ver:" + cfg.Apps[b] + ",Remote Ver:" + v + ")..")
								go appupdate(b,v)
							}
						}else{
							seelog.Debug("[func:SendToHeartBeat] App '"+ b +"' Need Install(Ver:"+v+")..")
							go appupdate(b,v)
						}
					}
					for localk,localv := range cfg.Apps{
						if _, ok := rdat.Data[localk]; !ok {
							seelog.Debug("[func:SendToHeartBeat] App '"+ localk +"' Will Remove(Ver:"+localv+")..")
							delete(cfg.Apps,localk)
							g.SaveConfig()
						}
					}
				}else{
					seelog.Error("[func:SendToHeartBeat] ",err)
					go g.SaveErrorStat("[func:SendToHeartBeat] "+err.Error())
				}
			}else{
				seelog.Error("[func:SendToHeartBeat] ",err)
				go g.SaveErrorStat("[func:SendToHeartBeat] "+err.Error())
			}
		}
	}
	seelog.Info("Send Info To HeartBeat Finish...")
}

func appupdate(appname string ,appver string){
	seelog.Info("App["+appname+"] Update Start..")
	seelog.Debug("[func:appupdate] Apps Download..("+g.Config().AppsDownloadAddr+"/"+appname+"_"+appver+".zip"+")")
	if err:=os.RemoveAll(g.Root+"/apps/"+appname);err==nil{
		if err:=funcs.Download(g.Config().AppsDownloadAddr+"/"+appname+"_"+appver+".zip",g.Root+"/apps/"+appname,"app-worker-package.zip");err==nil{
			if err:=funcs.Unzip(g.Root+"/apps/"+appname+"/app-worker-package.zip",g.Root+"/apps/"+appname);err==nil{
				g.Config().Apps[appname] = appver
				g.SaveConfig()
			}else{
				seelog.Error("[func:appupdate] New Apps Unzip Error ",err)
				go g.SaveErrorStat("[func:appupdate] New Apps Unzip Error "+err.Error())
			}
		}else{
			seelog.Error("[func:appupdate] Apps ("+appname+") Download Error ",err)
			go g.SaveErrorStat("[func:appupdate]  Apps ("+appname+") Download Error "+err.Error())
		}
	}else{
		seelog.Error("[func:appupdate] Apps ("+appname+") Update Remove First Error ",err)
		go g.SaveErrorStat("[func:appupdate]  Apps ("+appname+") Update Remove First Error "+err.Error())
	}
	seelog.Info("App["+appname+"] Update Finish..")
}
