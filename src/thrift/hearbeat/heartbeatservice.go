// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package hearbeat

import (
	"bytes"
	"fmt"
	"github.com/gy-games-libs/go-thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type HeartbeatService interface {
	// Parameters:
	//  - Info
	HeartbeatPackage(info *AgentInfo) (r string, err error)
}

type HeartbeatServiceClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewHeartbeatServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *HeartbeatServiceClient {
	return &HeartbeatServiceClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewHeartbeatServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *HeartbeatServiceClient {
	return &HeartbeatServiceClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - Info
func (p *HeartbeatServiceClient) HeartbeatPackage(info *AgentInfo) (r string, err error) {
	if err = p.sendHeartbeatPackage(info); err != nil {
		return
	}
	return p.recvHeartbeatPackage()
}

func (p *HeartbeatServiceClient) sendHeartbeatPackage(info *AgentInfo) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("heartbeatPackage", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := HeartbeatServiceHeartbeatPackageArgs{
		Info: info,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *HeartbeatServiceClient) recvHeartbeatPackage() (value string, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "heartbeatPackage" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "heartbeatPackage failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "heartbeatPackage failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error1 error
		error1, err = error0.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error1
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "heartbeatPackage failed: invalid message type")
		return
	}
	result := HeartbeatServiceHeartbeatPackageResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type HeartbeatServiceProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      HeartbeatService
}

func (p *HeartbeatServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *HeartbeatServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *HeartbeatServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewHeartbeatServiceProcessor(handler HeartbeatService) *HeartbeatServiceProcessor {

	self2 := &HeartbeatServiceProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self2.processorMap["heartbeatPackage"] = &heartbeatServiceProcessorHeartbeatPackage{handler: handler}
	return self2
}

func (p *HeartbeatServiceProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x3.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x3

}

type heartbeatServiceProcessorHeartbeatPackage struct {
	handler HeartbeatService
}

func (p *heartbeatServiceProcessorHeartbeatPackage) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := HeartbeatServiceHeartbeatPackageArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("heartbeatPackage", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := HeartbeatServiceHeartbeatPackageResult{}
	var retval string
	var err2 error
	if retval, err2 = p.handler.HeartbeatPackage(args.Info); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing heartbeatPackage: "+err2.Error())
		oprot.WriteMessageBegin("heartbeatPackage", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = &retval
	}
	if err2 = oprot.WriteMessageBegin("heartbeatPackage", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Info
type HeartbeatServiceHeartbeatPackageArgs struct {
	Info *AgentInfo `thrift:"info,1" json:"info"`
}

func NewHeartbeatServiceHeartbeatPackageArgs() *HeartbeatServiceHeartbeatPackageArgs {
	return &HeartbeatServiceHeartbeatPackageArgs{}
}

var HeartbeatServiceHeartbeatPackageArgs_Info_DEFAULT *AgentInfo

func (p *HeartbeatServiceHeartbeatPackageArgs) GetInfo() *AgentInfo {
	if !p.IsSetInfo() {
		return HeartbeatServiceHeartbeatPackageArgs_Info_DEFAULT
	}
	return p.Info
}
func (p *HeartbeatServiceHeartbeatPackageArgs) IsSetInfo() bool {
	return p.Info != nil
}

func (p *HeartbeatServiceHeartbeatPackageArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *HeartbeatServiceHeartbeatPackageArgs) readField1(iprot thrift.TProtocol) error {
	p.Info = &AgentInfo{}
	if err := p.Info.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Info), err)
	}
	return nil
}

func (p *HeartbeatServiceHeartbeatPackageArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("heartbeatPackage_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *HeartbeatServiceHeartbeatPackageArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("info", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:info: ", p), err)
	}
	if err := p.Info.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Info), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:info: ", p), err)
	}
	return err
}

func (p *HeartbeatServiceHeartbeatPackageArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("HeartbeatServiceHeartbeatPackageArgs(%+v)", *p)
}

// Attributes:
//  - Success
type HeartbeatServiceHeartbeatPackageResult struct {
	Success *string `thrift:"success,0" json:"success,omitempty"`
}

func NewHeartbeatServiceHeartbeatPackageResult() *HeartbeatServiceHeartbeatPackageResult {
	return &HeartbeatServiceHeartbeatPackageResult{}
}

var HeartbeatServiceHeartbeatPackageResult_Success_DEFAULT string

func (p *HeartbeatServiceHeartbeatPackageResult) GetSuccess() string {
	if !p.IsSetSuccess() {
		return HeartbeatServiceHeartbeatPackageResult_Success_DEFAULT
	}
	return *p.Success
}
func (p *HeartbeatServiceHeartbeatPackageResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *HeartbeatServiceHeartbeatPackageResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *HeartbeatServiceHeartbeatPackageResult) readField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 0: ", err)
	} else {
		p.Success = &v
	}
	return nil
}

func (p *HeartbeatServiceHeartbeatPackageResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("heartbeatPackage_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *HeartbeatServiceHeartbeatPackageResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRING, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := oprot.WriteString(string(*p.Success)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.success (0) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *HeartbeatServiceHeartbeatPackageResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("HeartbeatServiceHeartbeatPackageResult(%+v)", *p)
}
