package http

import (
	"crypto/md5"
	"fmt"
	"github.com/elves-project/agent/src/funcs"
	"github.com/elves-project/agent/src/g"
	"github.com/elves-project/agent/src/thrift/scheduler"
	"github.com/gy-games-libs/seelog"
	"github.com/gy-games-libs/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func configApiRoutes() {

	http.HandleFunc("/api/v2/rt/exec", func(w http.ResponseWriter, r *http.Request) {
		if g.Config().Devmode.Enabled {
			r.ParseForm()
			ret := g.ApiErrresut{}
			ret.Flag = "false"
			rert := g.ApiExecresult{}
			ins := scheduler.Instruct{}
			pastr := map[string]string{}
			if len(r.Form["auth_id"]) <= 0 {
				ret.Error = "[401.2]Unauthorized AuthId Is Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["sign"]) <= 0 || len(r.Form["sign"][0]) != 32 {
				ret.Error = "[401.4]Unauthorized Sign Length Is Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["sign_type"]) <= 0 || strings.ToLower(r.Form["sign_type"][0]) != "md5" {
				ret.Error = "[401.3]Unauthorized SignType Is Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["timestamp"]) <= 0 {
				ret.Error = "[401.8]Unauthorized Timestamp Is Illegal"
				RenderJson(w, ret)
				return
			}
			if _, err := strconv.Atoi(r.Form["timestamp"][0]); err != nil {
				ret.Error = "[401.8]Unauthorized Timestamp Is Illegal"
				RenderJson(w, ret)
				return
			}
			if funcs.StrToInt(r.Form["timestamp"][0]) < funcs.StrToInt(strconv.FormatInt(time.Now().Unix(), 10))-60*3 || funcs.StrToInt(r.Form["timestamp"][0]) > funcs.StrToInt(strconv.FormatInt(time.Now().Unix(), 10))+60*3 {
				ret.Error = "[401.7]Unauthorized Sign Timeout"
				RenderJson(w, ret)
				return
			}
			if r.Form["auth_id"][0] != g.Config().Devmode.Authid || len(r.Form["ip"]) > 0 && r.Form["ip"][0] != g.Config().Ip {
				ret.Error = "[401.1]Request Params (app) Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["ip"]) <= 0 {
				ret.Error = "[403.1]Request Params (ip) Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["app"]) <= 0 {
				ret.Error = "[403.3]Request Params (app) Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["func"]) <= 0 {
				ret.Error = "[403.7]Request Params (func) Illegal"
				RenderJson(w, ret)
				return
			}
			if len(r.Form["timeout"]) > 0 {
				if _, err := strconv.Atoi(r.Form["timeout"][0]); err != nil && r.Form["timeout"][0] != "" {
					ret.Error = "[403.8]Request Params (timeout) Illegal"
					RenderJson(w, ret)
					return
				}
			}
			pastr["app"] = r.Form["app"][0]
			pastr["func"] = r.Form["func"][0]
			if len(r.Form["param"]) > 0 {
				pastr["param"] = r.Form["param"][0]
				ins.Param = r.Form["param"][0]
			}
			if len(r.Form["timeout"]) > 0 {
				pastr["timeout"] = r.Form["timeout"][0]
				ins.Timeout = int32(funcs.StrToInt(pastr["timeout"]))
			}
			if len(r.Form["proxy"]) > 0 {
				pastr["proxy"] = r.Form["proxy"][0]
				ins.Proxy = r.Form["proxy"][0]
			}
			pastr["ip"] = g.Config().Ip
			pastr["auth_id"] = g.Config().Devmode.Authid
			pastr["timestamp"] = r.Form["timestamp"][0]
			sign, str := funcs.Sign("/api/v2/rt/exec?", pastr, g.Config().Devmode.Authkey)
			seelog.Debug("[funcs:/api/rt/create] Local ", str, " ", sign)
			seelog.Debug("[funcs:/api/rt/create] Get ", sign)
			if sign != r.Form["sign"][0] {
				ret.Error = "[401.5]Unauthorized Sign Error"
				RenderJson(w, ret)
				return
			}
			mret := g.Apiresut{}
			ins.ID = fmt.Sprintf("%x", md5.Sum([]byte(uuid.Rand().Hex())))[0:16]
			ins.IP = g.Config().Ip
			ins.Type = "rt"
			ins.App = r.Form["app"][0]
			ins.Func = r.Form["func"][0]
			f, r, c := funcs.Appexec(ins)
			rert.Rt_id = ins.ID
			rert.Worker_message = r
			rert.Worker_costtime = c
			rert.Worker_flag = strconv.Itoa(int(f))
			mret.Flag = "true"
			mret.Error = ""
			mret.Result = rert
			RenderJson(w, mret)
		}
	})

	http.HandleFunc("/api/gettesturl", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ret := map[string]string{}
		if g.Config().Devmode.Enabled {
			ret["status"] = "true"
			pastr := map[string]string{}
			if len(r.Form["app"]) > 0 && len(r.Form["func"]) > 0 {
				pastr["app"] = r.Form["app"][0]
				pastr["func"] = r.Form["func"][0]
				if len(r.Form["param"]) > 0 {
					pastr["param"] = r.Form["param"][0]
				}
				if len(r.Form["timeout"]) > 0 {
					pastr["timeout"] = r.Form["timeout"][0]
				}
				if len(r.Form["proxy"]) > 0 {
					pastr["proxy"] = r.Form["proxy"][0]
				}
				pastr["ip"] = g.Config().Ip
				pastr["auth_id"] = g.Config().Devmode.Authid
				pastr["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
				sign, str := funcs.Sign("/api/v2/rt/exec?", pastr, g.Config().Devmode.Authkey)
				ret["uri"] = str + "&sign_type=MD5&sign=" + sign
			} else {
				ret["uri"] = "param error!"
			}
		} else {
			ret["status"] = "false"
		}

		RenderDataJson(w, ret)
	})

}
