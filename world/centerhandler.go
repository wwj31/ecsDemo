package world

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message"
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/internal/message"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"ecsDemo/world/utils"
)

// world 通知center 实体切换区域
func (s *CenterWorld) W2CenterUpdateEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2CenterUpdateEntity)
	s.entities[data.EId] = sourceId
	logger.KVs(log.Fields{"eid": data.EId, "sourceId": sourceId, "areaId": data.AreaId, "msgName": tools.MsgName(data)}).Debug("world 通知center 实体切换区域")
}

// world 通知center 删除实体
func (s *CenterWorld) W2CenterDeleteEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2CenterDeleteEntity)
	delete(s.entities, data.EId)
	logger.KVs(log.Fields{"actorId": s.GetID(), "eid": data.EId, "sourceId": sourceId, "msgName": tools.MsgName(data)}).Debug("W2CenterDeleteEntity")
}

// game 通知center 玩家创建世界
func (s *CenterWorld) G2WCreateNewPlayer(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.G2WCreateNewPlayer)
	i, _, _ := utils.GridIndex(tools.Vec3f{X: data.X, Y: data.Y})
	s.player[data.GateSession] = data.RID
	s.Send(common.WorldName(i), data) // 通知world创建实体
	logger.KVs(log.Fields{"eid": data.EID, "sourceId": sourceId, "areaId": i, "msgName": tools.MsgName(data)}).Debug("G2WCreateNewPlayer")
}

// game 通知center 玩家进入世界
func (s *CenterWorld) G2WEnterPlayer(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.G2WEnterPlayer)
	s.player[data.GateSession] = data.RID
	entactor, ok := s.entities[data.EID]
	logf := log.Fields{"eid": data.EID, "sourceId": sourceId, "areaId": entactor, "msgName": tools.MsgName(data)}
	expect.True(ok, logf)

	s.Send(entactor, data) // 通知world玩家进入游戏
	logger.KVs(logf).Debug("G2WEnterPlayer")

}

// game 通知center 玩家进入世界
func (s *CenterWorld) G2WInvaildSession(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.G2WInvaildSession)
	delete(s.player, data.GateSession)
	logger.KVs(log.Fields{"session": data.GateSession}).Debug("G2WInvaildSession")
}

//////////////////////////////////////////////// 客户端消息 /////////////////////////////////////////////////////////////
// 操作实体移动
func (s *CenterWorld) WorldOperateEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*message.WorldOperateEntity)
	atrId, ok := s.entities[data.EID]
	expect.True(ok, log.Fields{"eid": data.EID})
	s.Send(atrId, inner_message.NewGateWrapperByPb(data, gateSession))
	logger.KVs(log.Fields{"eid": data.EID}).Green().Debug("centerworld receive operate")
}
