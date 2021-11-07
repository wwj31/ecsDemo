package entity

import (
	"github.com/vmihailenco/msgpack"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/internal/message"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

// 副本数据
func AreaDuplicate(w interfaces.IWorld, eid string) *inner.EntityInfo {
	var (
		data      *inner.EntityInfo
		comps     []ecs.ComponentType
		serialize map[uint64][]byte
		err       error
		ent       *ecs.Entity
	)
	serialize = make(map[uint64][]byte, 0)
	ent = w.Runtime().GetEntity(eid)
	expect.True(ent != nil, log.Fields{"eid": eid, "world": w.GetID()})

	data = &inner.EntityInfo{EId: ent.Id()}
	comps = []ecs.ComponentType{
		component.MOVE_COMP,
		component.AREA_COMP,
		component.POS_COMP,
		component.ATTRIBUTE_COMP,
		component.PLAY_COMP,
		component.FIGHTING_COMP,
	}

	f := func(t ecs.ComponentType, c ecs.IComponent) {
		serialize[t.ComponentType()], err = msgpack.Marshal(c)
		expect.Nil(err, log.Fields{"t": t})
	}

	ent.RangeComponent(f, comps...)
	data.Data = serialize
	return data
}

// 客户端实体全量所需数据
var TotalSyncData = []ecs.ComponentType{component.MOVE_COMP, component.POS_COMP, component.ATTRIBUTE_COMP}

// 打包实体消息
func EntityMsg(ent *ecs.Entity) *message.EntityData {
	entMsg := &message.EntityData{
		EID: ent.Id(),
	}
	for _, c := range TotalSyncData {
		comp, ok := ent.GetComponent(c).(interfaces.IMsgComponent)
		if ok {
			entMsg.Property = append(entMsg.Property, comp.SyncData())
		}
	}
	return entMsg
}
