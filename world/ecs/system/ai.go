package system

import (
	"ecsDemo/world/constant"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/interfaces"
	"ecsDemo/world/utils"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"reflect"
)

type AITuple struct {
	//fightComp *component.Fighting
	posComp   *component.Position
	areaComp  *component.AreaInfo
	attriComp *component.Attribute
}

func (s *AITuple) Init(comps map[ecs.ComponentType]ecs.IComponent) {
	var ok bool
	s.areaComp, ok = comps[s.areaComp.Type()].(*component.AreaInfo)
	expect.True(ok)

	s.posComp, ok = comps[s.posComp.Type()].(*component.Position)
	expect.True(ok)
	s.attriComp, ok = comps[s.attriComp.Type()].(*component.Attribute)
	expect.True(ok)
}

// 实体输入系统
type AISys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewAISys(w interfaces.IWorld) ecs.ISystem {
	ns := &AISys{
		world:      w,
		SystemBase: ecs.NewSystemBase(CLEAR_SYSTEM, w.Runtime(), reflect.TypeOf((*AITuple)(nil))),
	}
	return ns
}

func (s *AISys) EssentialComp() uint64 {
	return component.POS_COMP.ComponentType() | component.AREA_COMP.ComponentType()
}

func (s *AISys) UpdateFrame(float64) {
	s.Range(func(eid1 string, t1 ecs.ITuple) bool {
		tuple1, ok := t1.(*AITuple)
		expect.True(ok, log.Fields{"eid": eid1})
		if tuple1.areaComp.AreaId != s.world.AreaId() {
			return true
		}
		if tuple1.attriComp.Target != "" {
			return true
		}

		// 给eid1 找目标
		pos1 := tuple1.posComp.Pos
		s.Range(func(eid2 string, t2 ecs.ITuple) bool {
			if eid1 == eid2 {
				return true
			}
			tuple2, ok := t2.(*AITuple)
			expect.True(ok, log.Fields{"eid": eid1})
			pos2 := tuple2.posComp.Pos
			dist := tools.Distance(pos1, pos2)
			if dist < constant.ATTACK_DIST {
				utils.GetDataWithBuild(s.world, eid1).Attack = &component.AttackInputData{TargetEId: eid2}
				// ai找到目标后,如果没有进战斗就初始化战斗组件
				ent1 := s.world.Runtime().GetEntity(eid1)
				if ent1.GetComponent(component.FIGHTING_COMP) == nil {
					fightComp := component.NewFighting(NewFightState(s.world, eid1), utils.NewTurnCount())
					ent1.SetComponent(fightComp)
				}
				return false
			}
			return true
		})
		return true
	})
}
