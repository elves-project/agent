package g

import (
	"../thrift/scheduler"
	"sync"
	"time"
)

type StatGInfo struct {
	Mode   string
	Asset  string
	Ip     string
	Uptime string
	Hbtime string
	Ver    string
	Apps   map[string]string
}

type sIns struct {
	Time     string
	ID       string
	Type     string
	Mode     string
	Proxy    string
	App      string
	Func     string
	Costtime int32
	Flag     string
}

type StatCInfo struct {
	Id       string
	App      string
	Func     string
	Mode     string
	Rule     string
	Comment  string
	Lastexec string
}

var (
	stat  = StatGInfo{}
	tstat = [40]sIns{}
	estat = [20]string{}
	cstat = map[string]StatCInfo{}
	slock = new(sync.RWMutex)
	tlock = new(sync.RWMutex)
	elock = new(sync.RWMutex)
	clock = new(sync.RWMutex)
)

func InitStat() {
	s := StatGInfo{}
	s.Mode = "PRODUCT"
	if Config().Devmode.Enabled == true {
		s.Mode = "DEVELOP"
	}
	s.Asset = Config().Asset
	s.Ip = Config().Ip
	tm := time.Unix(time.Now().Unix(), 0)
	s.Uptime = tm.Format("2006/01/02 15:04:05")
	s.Ver = VERSION
	stat = s
}

func UpdateHbTime() {
	slock.Lock()
	defer slock.Unlock()
	tm := time.Unix(time.Now().Unix(), 0)
	stat.Hbtime = tm.Format("2006/01/02 15:04:05")
}

func SaveTaskStat(ins scheduler.Instruct, flag int32, costtime int32) {
	tlock.Lock()
	defer tlock.Unlock()
	tptask := tstat
	tm := time.Unix(time.Now().Unix(), 0)
	tstat[0].Time = tm.Format("2006/01/02 15:04:05")
	tstat[0].App = ins.App
	tstat[0].Func = ins.Func
	tstat[0].ID = ins.ID
	tstat[0].Mode = ins.Mode
	tstat[0].Type = ins.Type
	tstat[0].Proxy = ins.Proxy
	if flag == 0 {
		tstat[0].Flag = "failure"
	} else if flag == -1 {
		tstat[0].Flag = "error"
	} else {
		tstat[0].Flag = "success"
	}
	tstat[0].Costtime = costtime
	for i := 1; i < 40; i++ {
		tstat[i] = tptask[i-1]
	}
}

func SaveErrorStat(err string) {
	elock.Lock()
	defer elock.Unlock()
	errlist := estat
	tm := time.Unix(time.Now().Unix(), 0)
	estat[0] = tm.Format("2006/01/02 15:04:05") + "|" + err
	for i := 1; i < 20; i++ {
		estat[i] = errlist[i-1]
	}

}

func UpdateCronSata(id string) {
	clock.Lock()
	defer clock.Unlock()
	tm := time.Unix(time.Now().Unix(), 0)
	timestr := tm.Format("2006/01/02 15:04:05")
	SCI := cstat[id]
	//SCI1 := SCI[id]
	SCI.Lastexec = timestr
	cstat[id] = SCI
}

func SaveCronSata(id string, app string, funcs string, mode string, rule string, comment string) {
	clock.Lock()
	defer clock.Unlock()
	timestr := "0000/00/00 00:00:00"
	cstat[id] = StatCInfo{
		id,
		app,
		funcs,
		mode,
		rule,
		comment,
		timestr,
	}
}

func DelCronSata(id string) {
	//cstat = map[string]string{}
	clock.Lock()
	defer clock.Unlock()
	delete(cstat, id)
}

func GetCStat() map[string]StatCInfo {
	clock.Lock()
	defer clock.Unlock()
	return cstat
}

func GetAStat() map[string]string {
	return Config().Apps
}

func GetGStat() StatGInfo {
	slock.Lock()
	defer slock.Unlock()
	stat.Apps = Config().Apps
	return stat
}

func GetTStat() [40]sIns {
	tlock.Lock()
	defer tlock.Unlock()
	return tstat
}

func GetEStat() [20]string {
	elock.Lock()
	defer elock.Unlock()
	return estat
}
