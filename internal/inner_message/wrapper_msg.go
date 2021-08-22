package inner_message

import (
	"ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func NewGateWrapperByPb(pb proto.Message, gateSession string) *inner.GateMsgWrapper {
	data, err := proto.Marshal(pb)
	if err != nil {
		log.KV("err", err).ErrorStack(3, "marshal pb failed")
		return nil
	}
	return &inner.GateMsgWrapper{GateSession: gateSession, MsgName: tools.MsgName(pb), Data: data}
}

func NewGateWrapperByBytes(data []byte, msgName, gateSession string) *inner.GateMsgWrapper {
	return &inner.GateMsgWrapper{GateSession: gateSession, MsgName: msgName, Data: data}
}

func UnwrapperGateMsg(msg interface{}) (interface{}, string, error) {
	wrapper, is := msg.(*inner.GateMsgWrapper)
	if !is {
		return msg, "", nil
	}

	tp, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(wrapper.MsgName))
	if err != nil {
		log.KV("msgName", wrapper.MsgName).KV("error", err).Error("not find")
		return nil, wrapper.GateSession, err
	}

	actMsg := tp.New().Interface().(proto.Message)
	err = proto.Unmarshal(wrapper.Data, actMsg.(proto.Message))
	if err != nil {
		log.KV("MsgName", wrapper.MsgName).KV("error", err).KV("err", err).Error("Unmarshal failed")
		return nil, wrapper.GateSession, err
	}
	return actMsg, wrapper.GateSession, nil
}
