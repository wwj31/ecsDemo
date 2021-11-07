package utils

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/ecsDemo/internal/message"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

/*
	如果同步删除实体给前端，不传cs
*/
func SyncData(world interfaces.IWorld, eid string, syncGrids []int32, cs ...ecs.ComponentType) {
	sync_comp := world.Runtime().SingleComponent(component.SYNCINFO_COMP).(*component.SyncInfo)
	syncData := component.NewSyncData()
	syncData.EId = eid
	syncData.SyncGrids = syncGrids
	for _, c := range cs {
		syncData.Comp[c] = true
	}
	sync_comp.Data = append(sync_comp.Data, syncData)
	//log.KVs(log.Fields{"fc": world.FC(), "actorId": world.GetID(), "eid": eid}).Debug("syncData")
}

func PbData(world interfaces.IWorld, eid string, comps map[ecs.ComponentType]bool, EntityMsg *message.EntityData) {
	ent := world.Runtime().GetEntity(eid)
	for compType := range comps {
		msgComp, ok := ent.GetComponent(compType).(interfaces.IMsgComponent)
		expect.True(ok, log.Fields{"comps": compType})
		EntityMsg.Property = append(EntityMsg.Property, msgComp.SyncData())
	}
}
