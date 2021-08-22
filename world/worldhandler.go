package world

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message"
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/internal/message"
	"ecsDemo/internal/msgtools"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/ecs/entity"
	"ecsDemo/world/utils"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
)

// 创建新玩家
func (s *World) G2WCreateNewPlayer(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.G2WCreateNewPlayer)
	bornPos := tools.Vec3f{X: data.X, Y: data.Y}
	position := component.NewPosition(bornPos)
	move := component.NewMove(component.Pos(bornPos))
	areaId, _, _ := utils.GridIndex(bornPos)
	areaInfo := component.NewAreaInfo(areaId)
	playerComp := component.NewPlayer(data.RID, data.GateSession)
	attributeComp := component.NewAttribute()
	ent := ecs.NewEntity(data.EID)
	utils.SpawnEntity(s, ent, position, move, areaInfo, playerComp, attributeComp)
}

// 老玩家进入游戏
func (s *World) G2WEnterPlayer(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.G2WEnterPlayer)
	ent := s.Runtime().GetEntity(data.EID)
	expect.True(ent != nil, log.Fields{"eid": data.EID})
	playerComp := component.NewPlayer(data.RID, data.GateSession)
	ent.SetComponent(playerComp)
}

// 玩家视口
func (s *World) WorldWatchPosition(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*message.WorldWatchPosition)
	gridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	oldPos := utils.DelGridWatcher(gridComp, data.RID)
	watchPos := msgtools.Msg2Vec3f(data.Pos)
	areaComp := s.Runtime().SingleComponent(component.CROSS_AREA_COMP).(*component.CrossArea)
	x, y := utils.Pos(int(s.AreaId()))
	area := areaComp.Areas[x][y]
	if !utils.InBound(area.ExtendBound, watchPos) {
		//logger.KVs(log.Fields{"actorId": s.GetID(), "rid": data.RID, "watchPos": watchPos}).Debug("WorldWatchPosition")
		return
	}
	utils.AddGridWatcher(gridComp, watchPos, data.RID, gateSession)
	// 同步上一帧地块上所有实体数据给玩家
	updateMsg := &message.WorldUpdateEntity{}
	appears, disappears, _ := utils.SyncGrids(oldPos.X, oldPos.Y, watchPos.X, watchPos.Y)
	//grids := utils.AroundGrid(watchPos)

	syncgrids := []int32{}

	f := func(index int32, newOrDel bool) {
		grid, ok := gridComp.Grids[index]
		// 如果地块不归本区域管理，不需要同步
		if !ok {
			return
		}
		syncgrids = append(syncgrids, index)
		for _, ent := range grid.Entities {
			AreaComp, ok := ent.GetComponent(component.AREA_COMP).(*component.AreaInfo)
			if !ok {
				logger.Error("entity has not Area Component")
				continue
			}
			// 不是本区域的实体，不由本区域同步
			if AreaComp.AreaId != s.AreaId() {
				continue
			}
			if newOrDel {
				entMsg := entity.EntityMsg(ent)
				updateMsg.Entities = append(updateMsg.Entities, entMsg)
			} else {
				entMsg := &message.EntityData{EID: ent.Id()}
				updateMsg.Entities = append(updateMsg.Entities, entMsg)
			}
		}
	}

	for _, index := range appears {
		f(index, true)
	}
	for _, index := range disappears {
		f(index, false)
	}
	gateId, _ := common.SplitGateSession(gateSession)
	s.Send(gateId, inner_message.NewGateWrapperByPb(updateMsg, gateSession))
	logger.KVs(log.Fields{"actorId": s.GetID(), "rid": data.RID, "pos": watchPos, "syncgrids": syncgrids, "gateSession": gateSession}).Debug("player watch grid")
}

// 移动实体
func (s *World) WorldOperateEntity(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*message.WorldOperateEntity)
	ent := s.Runtime().GetEntity(data.EID)
	expect.True(ent != nil, log.Fields{"eid": ent.Id()})
	playerComp, ok := ent.GetComponent(component.PLAY_COMP).(*component.Player)
	expect.True(ok, log.Fields{"eid": ent.Id()})
	expect.True(playerComp.GateSession == gateSession, log.Fields{"eid": ent.Id(), "gateSession": gateSession})

	inputTuple := utils.GetDataWithBuild(s, data.EID)
	move := &component.MoveInputData{
		Path:  msgtools.Msg2Vec3f_Arr(data.Path),
		Speed: 30,
	}
	inputTuple.Move = move
	log.KVs(log.Fields{"eid": data.EID, "gateSession": gateSession}).Debug("player request move")
}
