package world

import (
	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/iniconfig"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/log/colorized"
	"github.com/wwj31/dogactor/tools"
	"github.com/golang/protobuf/proto"
	"github.com/wwj31/ecsDemo/internal/common"
	"github.com/wwj31/ecsDemo/internal/inner_message"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/internal/message"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

type (
	CenterWorld struct {
		actor.Base
		Config iniconfig.Config

		entities map[string]string // eid => actorId
		player   map[string]int64  // gateSession => RID
		msgMap   map[string]interfaces.RegistFun
	}
)

func (s *CenterWorld) OnInit() {
	s.entities = make(map[string]string)
	s.player = make(map[string]int64)
	s.msgMap = make(map[string]interfaces.RegistFun)

	// world  发来的消息
	s.RegistMsg((*inner.W2CenterUpdateEntity)(nil), s.W2CenterUpdateEntity)
	s.RegistMsg((*inner.W2CenterDeleteEntity)(nil), s.W2CenterDeleteEntity)

	// game  发来的消息
	s.RegistMsg((*inner.G2WCreateNewPlayer)(nil), s.G2WCreateNewPlayer)
	s.RegistMsg((*inner.G2WInvaildSession)(nil), s.G2WInvaildSession)
	s.RegistMsg((*inner.G2WEnterPlayer)(nil), s.G2WEnterPlayer)

	// client 发来的消息
	s.RegistMsg((*message.WorldOperateEntity)(nil), s.WorldOperateEntity)

	for i := int32(0); i < int32(9); i++ {
		s.System().Regist(actor.New(common.WorldName(i), &World{WorldId: i}, actor.SetMailBoxSize(10000),  actor.SetTimerAccuracy(100)))
	}
}

func (s *CenterWorld) RegistMsg(msg proto.Message, f interfaces.RegistFun) {
	msgName := tools.MsgName(msg)

	if _, ok := s.msgMap[msgName]; ok {
		log.KV("msgName", msgName).ErrorStack(3, "regist repeated message")
		return
	}
	s.msgMap[msgName] = f
}
func (s *CenterWorld) OnStop() bool {
	logger.KV("actor", s.GetID()).Info(colorized.Red("World stop!"))
	return true
}

func (s *CenterWorld) OnHandleMessage(sourceId, targetId string, msg interface{}) {
	actMsg, gateSession, _ := inner_message.UnwrapperGateMsg(msg)
	if pbMsg, ok := actMsg.(proto.Message); ok {
		msgName := tools.MsgName(pbMsg)
		handle := s.msgMap[msgName]
		expect.True(handle != nil, log.Fields{"msgName": msgName})
		handle(sourceId, actMsg, gateSession)
	}
}
