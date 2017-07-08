package g

type AppCfgStrct struct {
	Pro *ProStrct `json:"Processor"`
}

type ProStrct struct {
	Addr    string `json:"addr"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout"`
}

type ApiExecresult struct {
	Rt_id           string `json:"id"`
	Worker_message  string `json:"worker_message"`
	Worker_flag     string `json:"worker_flag"`
	Worker_costtime int32  `json:"worker_costtime"`
}

type Apiresut struct {
	Flag   string        `json:"flag"`
	Error  string        `json:"error"`
	Result ApiExecresult `json:"result"`
}

type ApiErrresut struct {
	Flag  string `json:"flag"`
	Error string `json:"error"`
}

type SchedulerInstruct struct {
	ID      string `thrift:"id,1" json:"id"`
	IP      string `thrift:"ip,2" json:"ip"`
	Type    string `thrift:"type,3" json:"type"`
	Mode    string `thrift:"mode,4" json:"mode"`
	App     string `thrift:"app,5" json:"app"`
	Func    string `thrift:"func,6" json:"func"`
	Param   string `thrift:"param,7" json:"param"`
	Timeout int32  `thrift:"timeout,8" json:"timeout"`
	Proxy   string `thrift:"proxy,9" json:"proxy"`
}
