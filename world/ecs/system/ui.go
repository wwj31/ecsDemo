package system

import (
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/internal/msgtools"
	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"reflect"
	"ecsDemo/world/constant"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/interfaces"
	"ecsDemo/world/utils"
)

type UITuple struct{}

func (s *UITuple) Init(comps map[ecs.ComponentType]ecs.IComponent) {}

type flag struct {
	*ecs.Entity
	bool
}

// ui系统
type UISys struct {
	*ecs.SystemBase
	world interfaces.IWorld
	ents  map[string]*flag
	up    bool
}

func NewUISys(w interfaces.IWorld) ecs.ISystem {
	ns := &UISys{
		world:      w,
		SystemBase: ecs.NewSystemBase(UI_SYSTEM, w.Runtime(), reflect.TypeOf((*UITuple)(nil))),
		ents:       make(map[string]*flag),
	}

	ns.world.RegistMsg((*inner.U2GAddEntity)(nil), ns.U2GAddEntity)
	ns.world.RegistMsg((*inner.U2GDelEntity)(nil), ns.U2GDelEntity)
	ns.world.RegistMsg((*inner.U2GMoveEntity)(nil), ns.U2GMoveEntity)

	ns.syncUI()
	return ns
}

func (s *UISys) syncUI() {
	s.up = true
	areaId := s.world.AreaId()
	x, y := utils.Pos(int(areaId))
	corssAreaComp := s.Runtime().SingleComponent(component.CROSS_AREA_COMP).(*component.CrossArea)
	GridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	area := corssAreaComp.Areas[x][y]
	points := []*inner.Vec3F{}
	for K, _ := range GridComp.Grids {
		p := msgtools.Vec3f2Inner(utils.DivGridKey(K))
		p.X = p.X * constant.GRID_SIZE
		p.Y = p.Y * constant.GRID_SIZE
		points = append(points, p)
	}
	msg := &inner.G2UArea{
		WorldId: int64(areaId),
		ExclusiveBound: &inner.Bound{
			Pos: &inner.Vec3F{
				X: area.ExclusiveBound.X,
				Y: area.ExclusiveBound.Y,
			},
			Width:  area.ExclusiveBound.Width,
			Height: area.ExclusiveBound.Height,
		},
		ActualBound: &inner.Bound{
			Pos: &inner.Vec3F{
				X: area.ActualBound.X,
				Y: area.ActualBound.Y,
			},
			Width:  area.ActualBound.Width,
			Height: area.ActualBound.Height,
		},
		Point: points,
	}
	s.world.Send("GraphUI", msg)
}
func (s *UISys) HandlerActorEvent(event interface{}) {
	switch evData := event.(type) {
	case *actor.Ev_newActor:
		if evData.ActorId == "GraphUI" {
			s.syncUI()
		}
	}
}

func (s *UISys) RemoveEntity(ent *ecs.Entity) {
	AreaInfo := ent.GetComponent(component.AREA_COMP).(*component.AreaInfo)
	msgReq := &inner.G2UDelEntity{
		EID:       ent.Id(),
		Duplicate: AreaInfo.AreaId != s.world.AreaId(),
		AreaId:    s.world.AreaId(),
	}
	s.world.Send("GraphUI", msgReq)
}
func (s *UISys) UpdateFrame(deltaTime float64) {
	if !s.up {
		return
	}
	count := 0
	s.Range(func(eid string, t ecs.ITuple) bool {
		ent := s.Runtime().GetEntity(eid)
		s.ents[eid] = &flag{
			Entity: ent,
			bool:   true,
		}
		AreaComp := ent.GetComponent(component.AREA_COMP).(*component.AreaInfo)
		posComp := ent.GetComponent(component.POS_COMP).(*component.Position)
		attriComp := ent.GetComponent(component.ATTRIBUTE_COMP).(*component.Attribute)
		msgReq := &inner.G2UUpdateEntity{
			EID:       eid,
			RealPos:   msgtools.Vec3f2Inner(posComp.Pos),
			AreaId:    s.world.AreaId(),
			Hp:        attriComp.HP(),
			Duplicate: s.world.AreaId() != AreaComp.AreaId,
		}
		s.world.Send("GraphUI", msgReq)
		count++
		//uiLog.KVs(log.Fields{"actorId": s.world.GetID(), "eid": eid, "area": AreaComp.AreaId, "posComp.Pos": posComp.Pos}).Debug("GraphUI system update")
		return true
	})

	msgReq := &inner.G2UCountEntity{
		WorldId: int64(s.world.AreaId()),
		Count:   int32(count),
	}
	s.world.Send("GraphUI", msgReq)
	for _, ent := range s.ents {
		if !ent.bool {
			s.RemoveEntity(ent.Entity)
		} else {
			ent.bool = false
		}
	}
}

func (s *UISys) EssentialComp() uint64 {
	return component.EVERYONES
}

// 添加实体
func (s *UISys) U2GAddEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.U2GAddEntity)
	bornPos := tools.Vec3f{X: data.RealPos.X, Y: data.RealPos.Y}
	// entity all component
	position := component.NewPosition(bornPos)
	move := component.NewMove(component.Pos(bornPos))
	mapinfo := component.NewAreaInfo(data.AreaNum)
	attri := component.NewAttribute()

	ent := ecs.NewEntity(data.EID)
	utils.SpawnEntity(s.world, ent, position, move, mapinfo, attri)
	//ent.SetComponent(position, move, mapinfo, attri)
	//err := s.Runtime().AddEntity(ent)
	//expect.Nil(err, log.Fields{"eid": ent.Id(), "err": err})
}

// 删除实体
func (s *UISys) U2GDelEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.U2GDelEntity)
	dead := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
	ent := s.Runtime().GetEntity(data.EID)
	if ent != nil {
		dead.Deads[data.EID] = ent
	}
	uiLog.Warn("U2GDelEntity")
}

// 处理网络输入操作
func (s *UISys) U2GMoveEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.U2GMoveEntity)
	entData := utils.GetDataWithBuild(s.world, data.EID)
	entData.Move = &component.MoveInputData{
		Path:  msgtools.Inner2Vec3f_Arr(data.Path),
		Speed: data.Speed,
	}
	uiLog.KVs(log.Fields{"world": s.world.GetID(), "eid": data.EID, "len": len(entData.Move.Path)}).Debug("U2GMoveEntity")
}
