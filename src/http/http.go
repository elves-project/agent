package http

import (
	"encoding/json"
	"fmt"
	"github.com/elves-project/agent/src/g"
	"github.com/gy-games-libs/seelog"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func init() {
	configPageRoutes()
	configStatRoutes()
	configApiRoutes()
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func Start() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("[funcs:http Start] ", err)
			go g.SaveErrorStat("[func:http Start] " + fmt.Sprintf("%s", err))
		}
	}()
	if !g.Config().Http.Enabled {
		seelog.Debug("Http Disable")
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		seelog.Debug("Http Addr Nil")
		return
	}

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	seelog.Info("elves agent web dashbord listening", addr)
	log.Fatalln(s.ListenAndServe())
}
