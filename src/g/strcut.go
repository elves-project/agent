package g

type AppCfgStrct struct{
	Pro 			*ProStrct		`json:"Processor"`
}

type ProStrct struct{
	Addr 			string			`json:"addr"`
	Port 			int			`json:"port"`
	Timeout 		int			`json:"timeout"`
}

type ApiExecresult struct {
	Rt_id 		string			`json:"id"`
	Worker_message 	string			`json:"worker_message"`
	Worker_flag 	string			`json:"worker_flag"`
	Worker_costtime int32			`json:"worker_costtime"`
}

type Apiresut struct {
	Flag   string			`json:"flag"`
	Error  string			`json:"error"`
	Result ApiExecresult		`json:"result"`
}

type ApiErrresut struct {
	Flag   string			`json:"flag"`
	Error  string			`json:"error"`
}
