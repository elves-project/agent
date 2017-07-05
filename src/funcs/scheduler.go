package funcs

import (
	"github.com/elves-project/agent/src/g"
	"github.com/elves-project/agent/src/thrift/scheduler"
	"github.com/gy-games-libs/go-thrift"
	"github.com/gy-games-libs/seelog"
	"strconv"
)

type AgentServiceServiceImpl struct{}

func (this *AgentServiceServiceImpl) InstructionInvokeAsync(insList []*scheduler.Instruct) (r []*scheduler.Reinstruct, err error) {

	for _, ins := range insList {
		reins := &scheduler.Reinstruct{}
		reins.Ins = ins
		go ExecAndSend(*ins, true)
		reins.Result_ = ""
		reins.Flag = 1
		r = append(r, reins)
	}
	return r, nil
}

func (this *AgentServiceServiceImpl) InstructionInvokeSync(ins *scheduler.Instruct) (r *scheduler.Reinstruct, err error) {
	seelog.Debug("[funcs:InstructionInvokeSync] ins:", ins)
	reins := &scheduler.Reinstruct{}
	ins.Type = "rt"
	ins.Mode = "np"
	reins.Ins = ins
	reins.Flag, reins.Result_, reins.Costtime = Appexec(*ins)
	return reins, nil
}

func (this *AgentServiceServiceImpl) AliveCheck() (r string, err error) {
	r = "success"
	return r, nil
}

func ServerRun() {

	transport, err := thrift.NewTServerSocket("0.0.0.0:" + strconv.Itoa(g.Config().Port))
	if err != nil {
		seelog.Error("[func:ServerRun] ", err)
		go g.SaveErrorStat("[func:ServerRun] " + err.Error())
	}

	transportFactory := thrift.NewTTransportFactory()
	if transportFactory == nil {
		seelog.Error("[func:ServerRun] ", "Failed to create new TransportFactory")
		go g.SaveErrorStat("[func:ServerRun] " + "Failed to create new TransportFactory")
		//return nil
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	hander := &AgentServiceServiceImpl{}
	processor := scheduler.NewAgentServiceProcessor(hander)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	seelog.Info("elves agent service listening:" + strconv.Itoa(g.Config().Port))
	server.Serve()

}
