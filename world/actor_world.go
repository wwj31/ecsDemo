package world

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message"
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/internal/message"
	"ecsDemo/world/constant"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/ecs/system"
	"ecsDemo/world/interfaces"
	"github.com/golang/protobuf/proto"
	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/log/colorized"
	"github.com/wwj31/dogactor/tools"
	"time"
)

type (
	World struct {
		actor.Base
		runtime       *ecs.Runtime
		fc            int64
		msgMap        map[string]interfaces.RegistFun
		worldStopping bool
		WorldId       int32
		timeset       map[string]int64 // 系统每帧执行时间 map[system]time
	}
)

func (s *World) OnInit() {
	s.timeset = make(map[string]int64)
	s.msgMap = make(map[string]interfaces.RegistFun)
	s.runtime = ecs.NewRuntime()
	// singleton component
	expect.Nil(s.runtime.AddSingleComponent(
		// 先初始化单件
		component.NewCrossArea(),     // 区间跨服相关
		component.NewSpawnEntity(),   // 处理创建对象
		component.NewDeadEntity(),    // 处理删除对象
		component.NewInputSet(),      // 收集实体输入
		component.NewMoverSet(),      // 收集移动的实体
		component.NewGrids(),         // 地图格子
		component.NewSyncGridsData(), //同步数据

	))

	// 添加顺序，决定系统执行顺序
	expect.Nil(s.runtime.AddSystem(
		// 优先更新已有实体逻辑
		system.NewSpawnSys(s), // 创建实体
		system.NewMoveSys(s),  // 移动系统
		system.NewAISys(s),
		system.NewFightingSys(s), // 战斗系统
		// todo ... 其他system

		// 处理实体，以及单例组件逻辑
		system.NewCrossAreaSys(s), // 区服边界处理系统
		system.NewUISys(s),        // 可视化UI
		system.NewGridSys(s),      // 地图地块系统  (必须Area之后)
		system.NewSyncSys(s),      // 同步数据给client
		system.NewDeadSys(s),      // 销毁实体
		system.NewClearSys(s),     // 清理

	))

	// 运行逻辑帧
	s.AddTimer(tools.UUID(),constant.FRAME_RATE*time.Millisecond, s.update,-1 )

	s.System().RegistEvent(s.GetID(), (*actor.Ev_newActor)(nil))

	s.RegistMsg((*inner.G2WCreateNewPlayer)(nil), s.G2WCreateNewPlayer)
	s.RegistMsg((*inner.G2WEnterPlayer)(nil), s.G2WEnterPlayer)
	s.RegistMsg((*message.WorldWatchPosition)(nil), s.WorldWatchPosition)
	s.RegistMsg((*message.WorldOperateEntity)(nil), s.WorldOperateEntity)
}
func (s *World) AreaId() int32 {
	return s.WorldId
}
func (s *World) Runtime() *ecs.Runtime {
	return s.runtime
}
func (s *World) update(dt int64) {
	s.fc++
	//log.KVs(log.Fields{"actorId": s.GetID(), "fc": s.fc}).Red().Debug("FRAME")
	rtime := s.runtime.Run(float64(dt)/float64(time.Millisecond), s.timeset)
	if rtime > constant.FRAME_RATE {
		log.KVs(log.Fields{"areaNum": s.AreaId(), "rtime": rtime, "timeset": s.timeset}).Warn("run time > 100")
	}
}

func (s *World) OnStop() bool {
	logger.KV("actor", s.GetID()).Info(colorized.Red("World stop!"))
	return true
}

func (s *World) FC() int64 {
	return s.fc
}

func (s *World) RegistMsg(msg proto.Message, f interfaces.RegistFun) {
	msgName := tools.MsgName(msg)
	if _, ok := s.msgMap[msgName]; ok {
		log.KV("msgName", msgName).ErrorStack(3, "regist repeated message")
		return
	}
	s.msgMap[msgName] = f
}

func (s *World) OnHandleMessage(sourceId, targetId string, msg interface{}) {
	if s.worldStopping && common.IsActorOf(common.Gate_Actor, sourceId) {
		return // 如果进入stopping状态，抛弃所有网关消息
	}

	actMsg, gateSession, _ := inner_message.UnwrapperGateMsg(msg)
	if pbMsg, ok := actMsg.(proto.Message); ok {
		msgName := tools.MsgName(pbMsg)

		handle := s.msgMap[msgName]
		expect.True(handle != nil, log.Fields{"msgName": msgName})
		handle(sourceId, actMsg, gateSession)
	}

}

func (s *World) OnHandleEvent(event interface{}) {
	s.Runtime().RangeSystem(func(iSystem ecs.ISystem) bool {
		if es, ok := iSystem.(interfaces.IEventSystem); ok {
			es.HandlerActorEvent(event)
		}
		return true
	})

}
