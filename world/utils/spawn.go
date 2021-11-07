package utils

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

func SpawnEntity(w interfaces.IWorld, new *ecs.Entity, comps ...ecs.IComponent) {
	for _, c := range comps {
		c.Init(new.Id())
	}
	spawnComp := w.Runtime().SingleComponent(component.SPAWN_COMP).(*component.Spawn)
	_, has := spawnComp.Newcomes[new.Id()]
	expect.True(has == false, log.Fields{"eid": new.Id()})
	spawnComp.Newcomes[new.Id()] = &component.NewEntity{Ent: new, Comps: comps}
	//log.KVs(log.Fields{"eid": new.Id()}).White().Debug("SpawnEntity")
}
