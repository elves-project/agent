package funcs

import (
	"../g"
	"../thrift/app"
	"../thrift/scheduler"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gy-games-libs/file"
	"github.com/gy-games-libs/go-thrift"
	"github.com/gy-games-libs/seelog"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ExecAndSend(ins scheduler.Instruct, t bool) {
	seelog.Info("App Ins Exec and Send Begin...")
	ret := scheduler.Reinstruct{}
	ret.Ins = &ins
	ret.Flag, ret.Result_, ret.Costtime = Appexec(ins)
	if t == true {
		go Resultsend(ret, g.Config().Scheduler.Addr, g.Config().Scheduler.Port, g.Config().Scheduler.Timeout)
	}
	if strings.ToLower(ins.Mode) == "p" {
		configContent, err := file.ToTrimString(g.Root + "/apps/" + ins.App + "/" + "appcfg.json")
		if err != nil {
			go g.SaveErrorStat("[Funcs:ExecAndSend] read app config file:" + g.Root + "/apps/" + ins.App + "/" + "appcfg.json" + "fail:" + err.Error())
			seelog.Error("[Funcs:ExecAndSend] read app config file:", g.Root+"/apps/"+ins.App+"/"+"appcfg.json", "fail:", err)
		}
		var c g.AppCfgStrct
		err = json.Unmarshal([]byte(configContent), &c)
		if err == nil {
			if c.Pro.Addr != "" && c.Pro.Port > 0 {
				go resultProcessorSend(ret, c.Pro.Addr, c.Pro.Port, c.Pro.Timeout)
			} else {
				go g.SaveErrorStat("[Funcs:ExecAndSend] pares " + ins.App + " app config file (Addr or Port) fail")
				seelog.Error("[Funcs:ExecAndSend] pares " + ins.App + " app config file (Addr or Port) fail")
			}
		} else {
			go g.SaveErrorStat("[Funcs:ExecAndSend] pares " + ins.App + " app config fail")
			seelog.Error("[Funcs:ExecAndSend] pares " + ins.App + " app config fail")
		}
	}
	seelog.Info("App Ins Exec and Send Finish...")
}

func Resultsend(ret scheduler.Reinstruct, addr string, port int, timeout int) {
	seelog.Info("App Exec Result Send Begin...")
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	seelog.Debug("[func:resultsend] Addr:" + addr + ":" + strconv.Itoa(port))
	transport, err := thrift.NewTSocket(net.JoinHostPort(addr, strconv.Itoa(port)))
	transport.SetTimeout(time.Duration(timeout) * time.Second)
	if err != nil {
		seelog.Error("[func:resultsend]", "Error resolving address:", err)
		go g.SaveErrorStat("[func:resultsend]" + "Error resolving address:" + err.Error())
	} else {
		client := scheduler.NewSchedulerServiceClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			seelog.Error("[func:resultsend]", "Error opening socket to "+addr+":"+string(port), " ", err)
			go g.SaveErrorStat("[func:resultsend]" + "Error opening socket to " + addr + ":" + string(port) + " " + err.Error())
		}
		defer transport.Close()
		schret, err := client.DataTransport(&ret)
		if err != nil {
			seelog.Error("[func:resultsend]App Exec and Send Result:", err)
			go g.SaveErrorStat("[func:resultsend]App Exec and Send Result:" + err.Error())
		} else {
			seelog.Debug("[func:resultsend]App Exec and Send Result:", schret)
		}
	}
	seelog.Info("App Exec Result Send Finish...")
}

func resultProcessorSend(ret scheduler.Reinstruct, addr string, port int, timeout int) {
	seelog.Info("App Exec Result Processor Send Begin...")
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	seelog.Debug("[func:resultProcessorSend] Addr:" + addr + ":" + strconv.Itoa(port))
	transport, err := thrift.NewTSocket(net.JoinHostPort(addr, strconv.Itoa(port)))
	transport.SetTimeout(time.Duration(timeout) * time.Second)
	if err != nil {
		seelog.Error("[func:resultProcessorSend]", "Error resolving address:", err)
		go g.SaveErrorStat("[func:resultProcessorSend]" + "Error resolving address:" + err.Error())
	} else {
		client := app.NewAppServiceClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			seelog.Error("[func:resultProcessorSend]", "Error opening socket to "+addr+":"+string(port), " ", err)
			go g.SaveErrorStat("[func:resultProcessorSend]" + "Error opening socket to " + addr + ":" + string(port) + " " + err.Error())
		}
		defer transport.Close()
		appins := app.Instruct{
			ID:      ret.Ins.ID,
			IP:      ret.Ins.IP,
			Type:    ret.Ins.Type,
			Mode:    ret.Ins.Mode,
			App:     ret.Ins.App,
			Func:    ret.Ins.Func,
			Param:   ret.Ins.Param,
			Timeout: ret.Ins.Timeout,
			Proxy:   ret.Ins.Proxy,
		}
		appret := app.Reinstruct{
			Ins:      &appins,
			Flag:     ret.GetFlag(),
			Costtime: ret.GetCosttime(),
			Result_:  ret.GetResult_(),
		}
		appproret, err := client.RunProcessor(&appret)
		if err != nil {
			seelog.Error("[func:resultProcessorSend]App Exec and Send Result:", err)
			go g.SaveErrorStat("[func:resultProcessorSend]App Exec and Send Result:" + err.Error())
		} else {
			seelog.Debug("[func:resultProcessorSend]App Exec and Send Result:", appproret)
		}
	}
	seelog.Info("App Exec Result Processor Send Finish...")
}

func Appexec(ins scheduler.Instruct) (int32, string, int32) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("[funcs:Appexec] ", err)
			go g.SaveErrorStat("[func:Appexec] " + fmt.Sprintf("%s", err))
		}
	}()
	ret := ""
	intflag := int32(1)
	seelog.Info("App(Id:" + ins.ID + ",Mode:" + ins.Mode + ",Type:" + ins.Type + ",App:" + ins.App + ",Func:" + ins.Func + ") Exec Start...")
	startTime := CurrentTimeMillis()
	timeoutflag := false
	appexi := false
	for k, _ := range g.Config().Apps {
		if k == ins.App {
			appexi = true
		}
	}
	if appexi != false {
		arg := []string{}
		execpath := ""
		if ins.Proxy != "" {
			pxy := strings.Split(ins.Proxy, "|")
			seelog.Debug("[funcs:Appexec] get proxy :", pxy)
			if len(pxy) == 2 {
				arg = append(arg, g.Root+"/apps/"+ins.App+"/"+pxy[1])
				arg = append(arg, ins.App)
				arg = append(arg, ins.Func)
				arg = append(arg, base64.StdEncoding.EncodeToString([]byte(ins.Param)))
				execpath = pxy[0]
			} else {
				arg = append(arg, ins.App)
				arg = append(arg, ins.Func)
				arg = append(arg, base64.StdEncoding.EncodeToString([]byte(ins.Param)))
				execpath = g.Root + "/apps/" + ins.App + "/" + pxy[0]
			}
		} else {
			switch os := runtime.GOOS; os {
			case "linux":
				//arg = append(arg, g.Root+"/bin/agentProxy.py")
				arg = append(arg, g.Root+"/apps/"+ins.App+"/app-worker.py")
				arg = append(arg, ins.App)
				arg = append(arg, ins.Func)
				arg = append(arg, base64.StdEncoding.EncodeToString([]byte(ins.Param)))
				execpath = "/usr/bin/python"
			default:
				arg = append(arg, ins.App)
				arg = append(arg, ins.Func)
				arg = append(arg, base64.StdEncoding.EncodeToString([]byte(ins.Param)))
				execpath = g.Root + "/apps/" + ins.App + "/app-worker.exe"
				//execpath = g.Root+"/bin/agentProxy.exe"
			}
		}
		seelog.Debug("[func:appexec] ", execpath, arg)
		cmd := exec.Command(execpath, arg...)
		seelog.Debug(ins.Timeout)
		if ins.Timeout > 0 {
			var timer *time.Timer
			timer = time.AfterFunc(time.Duration(ins.Timeout)*time.Second, func() {
				defer func() {
					if err := recover(); err != nil {
						seelog.Error("[funcs:Appexec AfterFunc] ", err)
						go g.SaveErrorStat("[func:Appexec AfterFunc] " + fmt.Sprintf("%s", err))
					}
				}()
				timer.Stop()
				cmd.Process.Kill()
				timeoutflag = true
			})
		}
		output, err := cmd.Output()
		if err != nil {
			seelog.Error("[func:appexec] ", ins.App, " ", ins.Func, " ", err)
			go g.SaveErrorStat("[func:appexec] " + ins.App + " " + ins.Func + " " + fmt.Sprintf("%s", err))
			intflag = -1
			ret = "app exec error : " + fmt.Sprintf("%s", err)
		} else {
			reg := regexp.MustCompile(`<ElvesWFlag>([\s\S]*)<\/ElvesWFlag>`)
			var result []byte
			if len(reg.FindSubmatch(output)) == 2 {
				flag := reg.FindSubmatch(output)[1]
				if BytesString(flag) == "true" {
					intflag = 1
				} else {
					intflag = 0
				}
				regc := regexp.MustCompile(`<ElvesWResult>([\s\S]*)<\/ElvesWResult>`)
				if len(regc.FindSubmatch(output)) == 2 {
					result = regc.FindSubmatch(output)[1]
					ret = BytesString(result)
					seelog.Debug("[func:appexec] ", ins, " rst: ", BytesString(result))
				} else {
					intflag = -1
					result = output
					ret = "app output result error : " + BytesString(output)
				}
			} else {
				intflag = -1
				ret = "app output flag error : " + BytesString(output)
			}
			ret = BytesString(result)
			if timeoutflag {
				intflag = -1
				ret = "app exec timeout"
			}
		}

	} else {
		seelog.Debug("[func:appexec] ", ins, " app not install")
		intflag = -1
		ret = "app [" + ins.App + "] not install !"
	}
	endTime := CurrentTimeMillis()
	seelog.Info("App(Id:" + ins.ID + ",Mode:" + ins.Mode + ",Type:" + ins.Type + ",App:" + ins.App + ",Func:" + ins.Func + ") Exec Finsh...")
	costtime := int32(endTime - startTime)
	go g.SaveTaskStat(ins, intflag, costtime)
	return intflag, ret, costtime
}
