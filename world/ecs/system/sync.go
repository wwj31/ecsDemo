package system

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message"
	"ecsDemo/internal/message"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/log"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/interfaces"
	"ecsDemo/world/utils"
)

//
type SyncSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewSyncSys(w interfaces.IWorld) ecs.ISystem {
	ns := &SyncSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(SYNC_SYSTEM, w.Runtime(), nil),
	}

	return ns
}
func (s *SyncSys) EssentialComp() uint64 {
	return component.FORBID_ENTITY
}
func (s *SyncSys) UpdateFrame(float64) {
	syncComp := s.Runtime().SingleComponent(component.SYNCINFO_COMP).(*component.SyncInfo)
	gridsComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)

	var sessionWithSync = make(map[string][]*message.EntityData)
	for _, v := range syncComp.Data {
		// 组装同步数据
		EntityMsg := &message.EntityData{EID: v.EId}
		utils.PbData(s.world, v.EId, v.Comp, EntityMsg)

		// 把同步数据加入所有Watcher集合中
		for _, gridkey := range v.SyncGrids {
			grid, ok := gridsComp.Grids[gridkey]
			if !ok {
				continue
			}
			for rid, _ := range grid.Watchers {
				watcher, ok := gridsComp.Watchers[rid]
				if !ok {
					continue
				}
				sessionWithSync[watcher.Session] = append(sessionWithSync[watcher.Session], EntityMsg)
			}
		}

	}

	for gateSession, data := range sessionWithSync {
		atrId, _ := common.SplitGateSession(gateSession)
		s.world.Send(atrId, inner_message.NewGateWrapperByPb(&message.WorldUpdateEntity{Entities: data}, gateSession))
		syncLog.KVs(log.Fields{"atrId": s.world.GetID(), "gateSession": gateSession, "entity len": len(data)}).Debug("sync data")
	}
}
