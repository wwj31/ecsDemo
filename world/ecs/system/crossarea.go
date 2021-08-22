package system

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"reflect"
	"ecsDemo/world/constant"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/ecs/entity"
	"ecsDemo/world/interfaces"
	"ecsDemo/world/utils"
)

type CrossTuple struct {
	posComp  *component.Position
	areaComp *component.AreaInfo
}

func (s *CrossTuple) Init(comps map[ecs.ComponentType]ecs.IComponent) {
	var ok bool
	s.areaComp, ok = comps[s.areaComp.Type()].(*component.AreaInfo)
	expect.True(ok)

	s.posComp, ok = comps[s.posComp.Type()].(*component.Position)
	expect.True(ok)
}

// 跨区域服系统
type CrossAreaSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewCrossAreaSys(w interfaces.IWorld) ecs.ISystem {
	ms := &CrossAreaSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(CROSS_AREA_SYSTEM, w.Runtime(), reflect.TypeOf((*CrossTuple)(nil))),
	}
	cross := w.Runtime().SingleComponent(component.CROSS_AREA_COMP).(*component.CrossArea)
	for i := 0; i < constant.SERVER_SPLIT_AREA*constant.SERVER_SPLIT_AREA; i++ {
		x := i % constant.SERVER_SPLIT_AREA
		y := i / constant.SERVER_SPLIT_AREA
		actb := component.Bound{
			X:      float64(x) * constant.SERVER_AREA_WIDTH,
			Y:      float64(y) * constant.SERVER_AREA_HEIGHT,
			Width:  constant.SERVER_AREA_WIDTH,
			Height: constant.SERVER_AREA_HEIGHT,
		}
		excb := component.Bound{
			X:      actb.X + constant.OVERLAP,
			Y:      actb.Y + constant.OVERLAP,
			Width:  constant.SERVER_AREA_WIDTH - constant.OVERLAP*2,
			Height: constant.SERVER_AREA_HEIGHT - constant.OVERLAP*2,
		}
		extb := component.Bound{
			X:      actb.X - constant.OVERLAP,
			Y:      actb.Y - constant.OVERLAP,
			Width:  constant.SERVER_AREA_WIDTH + constant.OVERLAP*2,
			Height: constant.SERVER_AREA_HEIGHT + constant.OVERLAP*2,
		}
		cross.Areas[x][y] = &component.Area{
			ExclusiveBound: excb,
			ActualBound:    actb,
			ExtendBound:    extb,
		}
	}
	ms.world.RegistMsg((*inner.W2WAddDuplicate)(nil), ms.W2WAddDuplicate)
	ms.world.RegistMsg((*inner.W2WDelDuplicate)(nil), ms.W2WDelDuplicate)
	ms.world.RegistMsg((*inner.W2WUpdateDuplicate)(nil), ms.W2WUpdateDuplicate)

	return ms
}

func (s *CrossAreaSys) UpdateFrame(deltaTime float64) {
	crossAreaComp := s.Runtime().SingleComponent(component.CROSS_AREA_COMP).(*component.CrossArea)
	inputComp := s.Runtime().SingleComponent(component.INPUT_COMP).(*component.InputSet)
	deadComp := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)

	s.Range(func(eid string, t ecs.ITuple) bool {
		tuple, ok := t.(*CrossTuple)
		expect.True(ok, log.Fields{"eid": eid})

		if tuple.posComp.Pos == tuple.posComp.OldPos && len(tuple.areaComp.DuplicatedAreas) > 0 {
			// 位置无变化，且已有副本区域，可以判断为是老实体并且没有位移，不需要计算
			return true
		}
		s.checkPos(eid, crossAreaComp, inputComp, tuple.areaComp, deadComp, tuple.posComp)
		return true
	})
}
func (s *CrossAreaSys) EssentialComp() uint64 {
	return component.POS_COMP.ComponentType() | component.AREA_COMP.ComponentType()
}

func (s *CrossAreaSys) checkPos(eid string,
	crossAreaComp *component.CrossArea,
	inputComp *component.InputSet,
	areaComp *component.AreaInfo,
	deadComp *component.DeadEntities,
	positionComp *component.Position) {

	localAreaId := s.world.AreaId()
	if areaComp.AreaId != localAreaId {
		return // 实体控制权再邻服，这里不需要做什么
	}

	var oldAreaId int32
	if positionComp.OldPos != tools.Invalid() {
		oldAreaId, _, _ = utils.GridIndex(positionComp.OldPos)
	}
	newAreaId, newx, newy := utils.GridIndex(positionComp.Pos)
	area := crossAreaComp.Areas[newx][newy]
	// 跨区域服处理(如果实体本帧已经死亡，不需要同步跨区域)
	if newAreaId != localAreaId && deadComp.Deads[eid] == nil {
		areaComp.AreaId = newAreaId
		neighborAreas := areaComp.DuplicatedAreas
		// 通知邻服，更新实体信息
		utils.SyncAreaDuplicate(entity.AreaDuplicate(s.world, eid), neighborAreas, s.world)
		// 通知中心服，实体所在区域
		msg := &inner.W2CenterUpdateEntity{EId: eid, AreaId: newAreaId}
		s.world.Send(common.Center_Actor, msg)
		crossLog.KVs(log.Fields{"neighbor": neighborAreas,
			"newPos":      positionComp.Pos,
			"newAreaId":   newAreaId,
			"localAreaId": localAreaId,
			"oldPos":      positionComp.OldPos,
			"eid":         eid,
		}).Yellow().Debug("cross area notify neighbor")

		return
	}

	if utils.InBound(area.ExclusiveBound, positionComp.OldPos) && inSharedArea(area, positionComp.Pos) { // 如果从ExclusiveBound走到ActualBound，通知邻服，创建实体副本
		neiphborIndex(area, int(localAreaId), positionComp.Pos, areaComp)
		if areaComp.DuplicatedAreas != nil {
			msg := &inner.W2WAddDuplicate{Entity: entity.AreaDuplicate(s.world, eid)}
			for n, _ := range areaComp.DuplicatedAreas {

				s.world.Send(common.WorldName(int32(n)), msg)
				crossLog.KVs(log.Fields{"eid": eid, "s": s.world.GetID(), "sendto": common.WorldName(int32(n)), "arr": areaComp.DuplicatedAreas}).Debug("广播添加实体")
			}
		}
	} else if utils.InBound(area.ExclusiveBound, positionComp.Pos) && inSharedArea(area, positionComp.OldPos) { // 如果从ActualBound走到ExclusiveBound，通知邻服，删除实体副本
		if areaComp.DuplicatedAreas != nil {
			msg := &inner.W2WDelDuplicate{EId: eid}
			for n, _ := range areaComp.DuplicatedAreas {
				s.world.Send(common.WorldName(int32(n)), msg)
				crossLog.KVs(log.Fields{"eid": eid, "s": s.world.GetID(), "r": common.WorldName(int32(n)), "arr": areaComp.DuplicatedAreas}).Debug("广播删除实体")
			}
		}
		areaComp.DuplicatedAreas = nil
	} else if (inSharedArea(area, positionComp.Pos) && inSharedArea(area, positionComp.OldPos)) ||
		oldAreaId != newAreaId { // 如果在ActualBound内移动,或者跨服了
		add, sub := neiphborIndex(area, int(localAreaId), positionComp.Pos, areaComp)
		if len(add) > 0 { // 通知新增的邻服，添加实体副本
			msg := &inner.W2WAddDuplicate{Entity: entity.AreaDuplicate(s.world, eid)}
			for _, id := range add {
				s.world.Send(common.WorldName(int32(id)), msg)
				crossLog.KVs(log.Fields{"eid": eid, "s": s.world.GetID(), "targetActor": common.WorldName(int32(id)), "add": add}).Debug("广播添加实体")
			}
		}
		if len(sub) > 0 { // 通知删除的邻服，删除实体副本
			msg := &inner.W2WDelDuplicate{EId: eid}
			for _, id := range sub {
				if int32(id) == s.world.AreaId() {
					continue
				}
				s.world.Send(common.WorldName(int32(id)), msg)
				crossLog.KVs(log.Fields{"eid": eid, "s": s.world.GetID(), "r": common.WorldName(int32(id)), "sub": sub}).Debug("广播删除实体")
			}
		}
	}
}

// 判断 pos 是否在 area 的共享区域
func inSharedArea(area *component.Area, pos tools.Vec3f) bool {
	if !utils.InBound(area.ExclusiveBound, pos) && utils.InBound(area.ActualBound, pos) {
		return true
	}
	return false
}

// area 判断的区域
// curIndex 当前区域号
// pos  如果pos在ExclusiveBound外，且在ActualBound内，计算需要共享实体的区域ID
// mapInfo 更新DuplicatedAreas
// add 新增的区域
// sub 删掉的区域
func neiphborIndex(area *component.Area, curIndex int, pos tools.Vec3f, mapInfo *component.AreaInfo) (add []int, sub []int) {
	if utils.InBound(area.ExclusiveBound, pos) || !utils.InBound(area.ActualBound, pos) {
		return
	}
	add = []int{}
	sub = []int{}
	eb := area.ExclusiveBound
	ab := area.ActualBound
	splitn := constant.SERVER_SPLIT_AREA

	temp := map[int]bool{}
	addfun := func(vals ...int) {
		for _, v := range vals {
			if v < 0 || v >= constant.SERVER_SPLIT_AREA*constant.SERVER_SPLIT_AREA {
				continue
			}
			if !mapInfo.DuplicatedAreas[v] {
				add = append(add, v)
			}
			temp[v] = true
		}
	}

	defer func() {
		for k, _ := range mapInfo.DuplicatedAreas {
			if !temp[k] {
				sub = append(sub, k)
			}
		}
		mapInfo.DuplicatedAreas = temp
	}()

	if eb.X <= pos.X && pos.X <= eb.X+eb.Width {
		if ab.Y <= pos.Y && pos.Y <= eb.Y {
			// 上
			addfun(curIndex - splitn)
		} else {
			// 下
			addfun(curIndex + splitn)
		}
		return
	}

	if eb.Y <= pos.Y && pos.Y <= eb.Y+eb.Height {
		if ab.X <= pos.X && pos.X <= eb.X {
			// 左
			if (curIndex % constant.SERVER_SPLIT_AREA) > 0 {
				addfun(curIndex - 1)
			}
		} else {
			// 右
			if (curIndex % constant.SERVER_SPLIT_AREA) < constant.SERVER_SPLIT_AREA-1 {
				addfun(curIndex + 1)
			}
		}
		return
	}

	if ab.X <= pos.X && pos.X <= eb.X {
		arr := make([]int, 0, 4)
		if (curIndex % constant.SERVER_SPLIT_AREA) > 0 {
			arr = append(arr, curIndex-1)
		}
		if ab.Y <= pos.Y && pos.Y <= eb.Y {
			// 左上
			arr = append(arr, curIndex-splitn, curIndex-splitn-1)
		} else {
			// 左下
			arr = append(arr, curIndex+splitn, curIndex+splitn-1)
		}
		addfun(arr...)
		return
	}

	if eb.X+eb.Width <= pos.X && pos.X <= ab.X+ab.Width {
		arr := make([]int, 0, 4)
		if (curIndex % constant.SERVER_SPLIT_AREA) < constant.SERVER_SPLIT_AREA-1 {
			arr = append(arr, curIndex+1)
		}
		if ab.Y <= pos.Y && pos.Y <= eb.Y {
			// 右上
			arr = append(arr, curIndex-splitn, curIndex-splitn+1)
		} else {
			// 右下
			arr = append(arr, curIndex+splitn, curIndex+splitn+1)
		}
		addfun(arr...)
	}
	return
}
