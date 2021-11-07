package system

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/world/constant"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/ecs/entity"
	"github.com/wwj31/ecsDemo/world/interfaces"
	"github.com/wwj31/ecsDemo/world/utils"
	"reflect"
)

type MoveTuple struct {
	moveComp *component.Move
	posComp  *component.Position
	areaComp *component.AreaInfo
}

func (s *MoveTuple) Init(comps map[ecs.ComponentType]ecs.IComponent) {
	var ok bool
	s.moveComp, ok = comps[s.moveComp.Type()].(*component.Move)
	expect.True(ok)

	s.posComp, ok = comps[s.posComp.Type()].(*component.Position)
	expect.True(ok)

	s.areaComp, ok = comps[s.areaComp.Type()].(*component.AreaInfo)
	expect.True(ok)
}

// 移动系统
type MoveSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewMoveSys(w interfaces.IWorld) ecs.ISystem {
	ms := &MoveSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(MOVE_SYSTEM, w.Runtime(), reflect.TypeOf((*MoveTuple)(nil))),
	}
	return ms
}

// deltaTime 时间增量
func (s *MoveSys) UpdateFrame(deltaTime float64) {
	inputComp := s.Runtime().SingleComponent(component.INPUT_COMP).(*component.InputSet)
	moverComp := s.Runtime().SingleComponent(component.MOVER_COMP).(*component.MoverSet)
	s.Range(func(eid string, t ecs.ITuple) bool {
		tuple, ok := t.(*MoveTuple)
		expect.True(ok, log.Fields{"eid": eid})

		tuple.posComp.OldPos = tuple.posComp.Pos
		if tuple.moveComp.Velocity > 0 && tuple.posComp.Pos != tuple.moveComp.DestPos {
			s.move(deltaTime, eid, tuple.moveComp, tuple.posComp, 0)
			moverComp.Movers = append(moverComp.Movers, eid)
		}

		// 路径改变，刷新新的路径
		if entData, ok := inputComp.Inputs[eid]; ok && entData.Move != nil {
			tuple.moveComp.Velocity = entData.Move.Speed
			tuple.moveComp.Path = entData.Move.Path
			s.flushPath(tuple.posComp, tuple.moveComp)
			utils.SyncData(s.world, eid, utils.AroundGrid(tuple.posComp.Pos), component.MOVE_COMP) // 改变移动

			utils.SyncAreaDuplicate(entity.AreaDuplicate(s.world, eid), tuple.areaComp.DuplicatedAreas, s.world)
		}
		return true
	})
	//	logger.KVs(logger.Fields{"stop": stop, "move": move, "s.world.AreaId()": s.world.AreaId()}).Warn("stop count")
}
func (s *MoveSys) EssentialComp() uint64 {
	return component.MOVE_COMP.ComponentType() | component.POS_COMP.ComponentType() | component.AREA_COMP.ComponentType()
}

////////////////////////////////////////////// 内部函数 /////////////////////////////////////////////////////////////
// 移动增量更新
// deltaTime   时间增量
// eid  实体id
// comp 移动组件
// n   递归安全计数
func (s *MoveSys) move(deltaTime float64, eid string, movecomp *component.Move, poscomp *component.Position, n int) {
	if deltaTime == 0 {
		return
	}
	if poscomp.Pos == movecomp.DestPos {
		return
	}

	// 两次帧率不一致，重新计算速度偏移值
	if deltaTime != movecomp.Lastf {
		moveLog.KVs(log.Fields{"eid": eid, "dt": deltaTime, "movecomp.Lastf": movecomp.Lastf}).Debug("deltaTime != movecomp.Lastf ")
		movecomp.Lastf = deltaTime
		calcOffset(movecomp, poscomp.Pos, deltaTime)
	}
	poscomp.Pos.Add(movecomp.Offset)
	moveLog.KVs(log.Fields{"areaid": s.world.AreaId(), "n": n, "eid": eid, "movecomp": movecomp, "poscomp": poscomp}).Debug("move")

	// 到达目的点后分帧处理:根据多走的路程(eLen),计算多走时间增量(eMs),根据eMs分帧处理新方向的移动
	if tlen := tools.Distance(poscomp.Pos, movecomp.OriginPos); tlen >= movecomp.Distance {
		poscomp.Pos = movecomp.DestPos
		eLen := tlen - movecomp.Distance
		s.flushPath(poscomp, movecomp)
		if eLen > 0 && movecomp.Velocity > 0 {
			sp := offest(movecomp.Velocity, 1) // 1毫秒的增量
			eMs := eLen / sp                   // 额外移动的时间增量
			moveLog.KVs(log.Fields{"sp": sp, "eMs": eMs}).Debug("split frame")
			s.move(eMs, eid, movecomp, poscomp, n+1)
		}
		moveLog.Debug("stop point")
	}
}

// 刷新新的直线移动
func (s *MoveSys) flushPath(pos *component.Position, move *component.Move) {
	move.OriginPos = pos.Pos
	if len(move.Path) > 0 {
		d := move.Path[0]
		move.Path = move.Path[1:]
		move.DestPos = d
	} else {
		move.DestPos = pos.Pos
		move.Velocity = 0
	}
	calcOffset(move, pos.Pos, constant.FRAME_RATE) // 新目的，刷新移动偏移量
}

// 根据帧率计算速度偏移量
func calcOffset(move *component.Move, curPos tools.Vec3f, f float64) {
	if move.DestPos == curPos || move.OriginPos == move.DestPos || move.Velocity == 0 {
		move.Offset.X = 0
		move.Offset.Y = 0
		move.Distance = 0
		return
	}

	v := tools.Sub(move.DestPos, curPos)
	offsetDis := offest(move.Velocity, f)
	move.Distance = tools.Distance(move.OriginPos, move.DestPos)
	move.Offset = tools.Mul(v, offsetDis/move.Distance)
}

// 速度转换成一帧偏移量
// v 速度
// f 时间增量(毫秒)
func offest(v, f float64) float64 {
	return (v / 1000) * f
}
