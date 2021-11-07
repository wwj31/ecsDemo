package utils

import (
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

const (
	EnterGrid = iota
	LeaveGrid
	MoveNW //西北方向
	MoveN  //北方向
	MoveNE //东北方向
	MoveW  //西方向
	SameGirds
	MoveE  //东方向
	MoveSW //西南向
	MoveS  //南方向
	MoveSE //东南方向
)

var _syncGridBench = map[int]*SyncGirds{
	EnterGrid: syncGridWithEnter(),
	LeaveGrid: syncLeaveGrid(),
	MoveNW:    moveNW(),
	MoveN:     moveN(),
	MoveNE:    moveNE(),
	MoveW:     moveW(),
	SameGirds: sameGirds(),
	MoveE:     moveE(),
	MoveSW:    moveSW(),
	MoveS:     moveS(),
	MoveSE:    moveSE(),
}

type SyncGirds struct {
	AppearArray    []*rowCol //出现
	DisappearArray []*rowCol //消失
	ModifyArr      []*rowCol //变化
}

type rowCol struct {
	changeRow int
	changeCol int
}

// 获取坐标所在9宫格地块
func AroundGrid(pos tools.Vec3f) []int32 {
	row, col := GetGridRowCol(pos.X, pos.Y)
	_, _, modifyArr := ActionWithMove(SameGirds, row, col)
	return modifyArr
}

func SyncGrids(fromX float64, fromY float64, toX float64, toY float64) (appearArray []int32, disappearArr []int32, modifyArr []int32) {
	fromRow, fromCol := GetGridRowCol(fromX, fromY)
	toRow, toCol := GetGridRowCol(toX, toY)
	var action int
	switch {
	case fromX == -1 && fromY == -1: //进入格子
		action = EnterGrid
		fromRow, fromCol = toRow, toCol
	case toX == -1 && toY == -1: //离开格子
		action = LeaveGrid
	case fromRow-1 == toRow && fromCol-1 == toCol: //西北方向
		action = MoveNW
	case fromRow-1 == toRow && fromCol == toCol: //北
		action = MoveN
	case fromRow-1 == toRow && fromCol+1 == toCol: //东北
		action = MoveNE
	case fromRow == toRow && fromCol-1 == toCol: //西方向
		action = MoveW
	case fromRow == toRow && fromCol == toCol:
		action = SameGirds
	case fromRow == toRow && fromCol+1 == toCol: //东方向
		action = MoveE
	case fromRow+1 == toRow && fromCol-1 == toCol: //西南向
		action = MoveSW
	case fromRow+1 == toRow && fromCol == toCol: //南方向
		action = MoveS
	case fromRow+1 == toRow && fromCol+1 == toCol: //东南方向
		action = MoveSE
	default:
		_, disappearArr, _ = ActionWithMove(LeaveGrid, fromRow, fromCol) //离开之前的，进入最新的
		appearArray, _, _ = ActionWithMove(EnterGrid, toRow, toCol)
		return
	}
	return ActionWithMove(action, fromRow, fromCol)
}

func ActionWithMove(action int, row int, col int) (appearArray []int32, disappearArr []int32, modifyArr []int32) {
	syncGrids, ok := _syncGridBench[action]
	if !ok {
		log.KVs(log.Fields{"action": action}).Error("no this action")
		return
	}
	f := func(rc *rowCol) (id int32, flag bool) {
		gridId := GetGridByRowCol(row+rc.changeRow, col+rc.changeCol)
		if gridId < 0 {
			return
		}
		return gridId, true
	}
	for _, change := range syncGrids.AppearArray {
		if gridId, ok := f(change); ok {
			appearArray = append(appearArray, gridId)
		}
	}
	for _, change := range syncGrids.DisappearArray {
		if gridId, ok := f(change); ok {
			disappearArr = append(disappearArr, gridId)
		}
	}
	for _, change := range syncGrids.ModifyArr {
		if gridId, ok := f(change); ok {
			modifyArr = append(modifyArr, gridId)
		}
	}
	return
}

//
/*
	NW:grid-col-1  row-1,col-1
	N:grid-col	   row-1,col
	NE:grid-col+1  row-1,col+1
	W:grid-1  	   row,col-1
	self :0    	   row,col
	E:grid+1	   row,col+1
	SW:grid+col-1  row+1,col-1
	S:grid+col	   row+1,col
	SE:grid+col+1  row+1,col+1

	NW (-1,-1)  N:北(-1,0) NE(-1,1)
	W:西(0,-1)  role(0,0)  E:东(0,1)
	SW (1,-1)   S:南(1,0)  SE(1,1)
*/

//直接给出变化数据，而不是差异比较移动后格子的相关变化，减少计算
//进入格子
func syncGridWithEnter() *SyncGirds {
	sync := &SyncGirds{}
	sync.AppearArray = make([]*rowCol, 9, 9)
	sync.AppearArray[0] = &rowCol{changeRow: -1, changeCol: -1} //NW
	sync.AppearArray[1] = &rowCol{changeRow: -1, changeCol: 0}  //N
	sync.AppearArray[2] = &rowCol{changeRow: -1, changeCol: 1}  ////NE
	sync.AppearArray[3] = &rowCol{changeRow: 0, changeCol: -1}  //W
	sync.AppearArray[4] = &rowCol{changeRow: 0, changeCol: 0}   //self
	sync.AppearArray[5] = &rowCol{changeRow: 0, changeCol: 1}   //E
	sync.AppearArray[6] = &rowCol{changeRow: 1, changeCol: -1}  //SW
	sync.AppearArray[7] = &rowCol{changeRow: 1, changeCol: 0}   //S
	sync.AppearArray[8] = &rowCol{changeRow: 1, changeCol: 1}   //SE
	return sync
}

//离开格子
func syncLeaveGrid() *SyncGirds {
	sync := &SyncGirds{}
	sync.DisappearArray = make([]*rowCol, 9, 9)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: -1} //NW
	sync.DisappearArray[1] = &rowCol{changeRow: -1, changeCol: 0}  //N
	sync.DisappearArray[2] = &rowCol{changeRow: -1, changeCol: 1}  ////NE
	sync.DisappearArray[3] = &rowCol{changeRow: 0, changeCol: -1}  //W
	sync.DisappearArray[4] = &rowCol{changeRow: 0, changeCol: 0}   //self
	sync.DisappearArray[5] = &rowCol{changeRow: 0, changeCol: 1}   //E
	sync.DisappearArray[6] = &rowCol{changeRow: 1, changeCol: -1}  //SW
	sync.DisappearArray[7] = &rowCol{changeRow: 1, changeCol: 0}   //S
	sync.DisappearArray[8] = &rowCol{changeRow: 1, changeCol: 1}   //SE
	return sync
}

// 移动到NW西北方向格子
func moveNW() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 5, 5)
	sync.AppearArray[0] = &rowCol{changeRow: -2, changeCol: -2}
	sync.AppearArray[1] = &rowCol{changeRow: -2, changeCol: -1}
	sync.AppearArray[2] = &rowCol{changeRow: -2, changeCol: 0}
	sync.AppearArray[3] = &rowCol{changeRow: -1, changeCol: -2}
	sync.AppearArray[4] = &rowCol{changeRow: 0, changeCol: -2}

	sync.DisappearArray = make([]*rowCol, 5, 5)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: 1}
	sync.DisappearArray[1] = &rowCol{changeRow: 0, changeCol: 1}
	sync.DisappearArray[2] = &rowCol{changeRow: 1, changeCol: 1}
	sync.DisappearArray[3] = &rowCol{changeRow: 1, changeCol: 0}
	sync.DisappearArray[4] = &rowCol{changeRow: 1, changeCol: -1}

	sync.ModifyArr = make([]*rowCol, 4, 4)
	sync.ModifyArr[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.ModifyArr[1] = &rowCol{changeRow: -1, changeCol: 0}
	sync.ModifyArr[2] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[3] = &rowCol{changeRow: 0, changeCol: -1}
	return sync
}

//N:北
// 移动到北方向格子
func moveN() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 3, 3)
	sync.AppearArray[0] = &rowCol{changeRow: -2, changeCol: -1}
	sync.AppearArray[1] = &rowCol{changeRow: -2, changeCol: 0}
	sync.AppearArray[2] = &rowCol{changeRow: -2, changeCol: 1}

	sync.DisappearArray = make([]*rowCol, 3, 3)
	sync.DisappearArray[0] = &rowCol{changeRow: 1, changeCol: -1}
	sync.DisappearArray[1] = &rowCol{changeRow: 1, changeCol: 0}
	sync.DisappearArray[2] = &rowCol{changeRow: 1, changeCol: 1}

	sync.ModifyArr = make([]*rowCol, 6, 6)
	sync.ModifyArr[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.ModifyArr[1] = &rowCol{changeRow: -1, changeCol: 0}
	sync.ModifyArr[2] = &rowCol{changeRow: -1, changeCol: 1}
	sync.ModifyArr[3] = &rowCol{changeRow: 0, changeCol: -1}
	sync.ModifyArr[4] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[5] = &rowCol{changeRow: 0, changeCol: 1}
	return sync
}

//NE:东北
// 移动到东北方向格子
func moveNE() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 5, 5)
	sync.AppearArray[0] = &rowCol{changeRow: -2, changeCol: 0}
	sync.AppearArray[1] = &rowCol{changeRow: -2, changeCol: 1}
	sync.AppearArray[2] = &rowCol{changeRow: -2, changeCol: 2}
	sync.AppearArray[3] = &rowCol{changeRow: -1, changeCol: 2}
	sync.AppearArray[4] = &rowCol{changeRow: 0, changeCol: 2}

	sync.DisappearArray = make([]*rowCol, 5, 5)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.DisappearArray[1] = &rowCol{changeRow: 0, changeCol: -1}
	sync.DisappearArray[2] = &rowCol{changeRow: 1, changeCol: -1}
	sync.DisappearArray[3] = &rowCol{changeRow: 1, changeCol: 0}
	sync.DisappearArray[4] = &rowCol{changeRow: 1, changeCol: 1}

	sync.ModifyArr = make([]*rowCol, 4, 4)
	sync.ModifyArr[0] = &rowCol{changeRow: -1, changeCol: 0}
	sync.ModifyArr[1] = &rowCol{changeRow: -1, changeCol: 1}
	sync.ModifyArr[2] = &rowCol{changeRow: 0, changeCol: 1}
	sync.ModifyArr[3] = &rowCol{changeRow: 0, changeCol: 0}
	return sync
}

//W:西
// 移动到西方向格子
func moveW() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 3, 3)
	sync.AppearArray[0] = &rowCol{changeRow: -1, changeCol: -2}
	sync.AppearArray[1] = &rowCol{changeRow: 0, changeCol: -2}
	sync.AppearArray[2] = &rowCol{changeRow: 1, changeCol: -2}

	sync.DisappearArray = make([]*rowCol, 3, 3)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: 1}
	sync.DisappearArray[1] = &rowCol{changeRow: 0, changeCol: 1}
	sync.DisappearArray[2] = &rowCol{changeRow: 1, changeCol: 1}

	sync.ModifyArr = make([]*rowCol, 6, 6)
	sync.ModifyArr[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.ModifyArr[1] = &rowCol{changeRow: -1, changeCol: 0}
	sync.ModifyArr[2] = &rowCol{changeRow: 0, changeCol: -1}
	sync.ModifyArr[3] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[4] = &rowCol{changeRow: 1, changeCol: -1}
	sync.ModifyArr[5] = &rowCol{changeRow: 1, changeCol: 0}
	return sync
}

//同格子移动
func sameGirds() *SyncGirds {
	sync := &SyncGirds{}
	sync.ModifyArr = make([]*rowCol, 9, 9)
	sync.ModifyArr[0] = &rowCol{changeRow: -1, changeCol: -1} //NW
	sync.ModifyArr[1] = &rowCol{changeRow: -1, changeCol: 0}  //N
	sync.ModifyArr[2] = &rowCol{changeRow: -1, changeCol: 1}  ////NE
	sync.ModifyArr[3] = &rowCol{changeRow: 0, changeCol: -1}  //W
	sync.ModifyArr[4] = &rowCol{changeRow: 0, changeCol: 0}   //self
	sync.ModifyArr[5] = &rowCol{changeRow: 0, changeCol: 1}   //E
	sync.ModifyArr[6] = &rowCol{changeRow: 1, changeCol: -1}  //SW
	sync.ModifyArr[7] = &rowCol{changeRow: 1, changeCol: 0}   //S
	sync.ModifyArr[8] = &rowCol{changeRow: 1, changeCol: 1}   //SE
	return sync
}

//E:东
// 移动到东方向格子
func moveE() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 3, 3)
	sync.AppearArray[0] = &rowCol{changeRow: -1, changeCol: 2}
	sync.AppearArray[1] = &rowCol{changeRow: 0, changeCol: 2}
	sync.AppearArray[2] = &rowCol{changeRow: 1, changeCol: 2}

	sync.DisappearArray = make([]*rowCol, 3, 3)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.DisappearArray[1] = &rowCol{changeRow: 0, changeCol: -1}
	sync.DisappearArray[2] = &rowCol{changeRow: 1, changeCol: -1}

	sync.ModifyArr = make([]*rowCol, 6, 6)
	sync.ModifyArr[0] = &rowCol{changeRow: -1, changeCol: 0}
	sync.ModifyArr[1] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[2] = &rowCol{changeRow: 1, changeCol: 0}
	sync.ModifyArr[3] = &rowCol{changeRow: -1, changeCol: 1}
	sync.ModifyArr[4] = &rowCol{changeRow: 0, changeCol: 1}
	sync.ModifyArr[5] = &rowCol{changeRow: 1, changeCol: 1}
	return sync
}

//SW:西南
// 移动到西南向格子
func moveSW() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 5, 5)
	sync.AppearArray[0] = &rowCol{changeRow: 0, changeCol: -2}
	sync.AppearArray[1] = &rowCol{changeRow: 1, changeCol: -2}
	sync.AppearArray[2] = &rowCol{changeRow: 2, changeCol: -2}
	sync.AppearArray[3] = &rowCol{changeRow: 2, changeCol: -1}
	sync.AppearArray[4] = &rowCol{changeRow: 2, changeCol: 0}

	sync.DisappearArray = make([]*rowCol, 5, 5)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.DisappearArray[1] = &rowCol{changeRow: -1, changeCol: 0}
	sync.DisappearArray[2] = &rowCol{changeRow: -1, changeCol: 1}
	sync.DisappearArray[3] = &rowCol{changeRow: 0, changeCol: 1}
	sync.DisappearArray[4] = &rowCol{changeRow: 1, changeCol: 1}

	sync.ModifyArr = make([]*rowCol, 4, 4)
	sync.ModifyArr[0] = &rowCol{changeRow: 0, changeCol: -1}
	sync.ModifyArr[1] = &rowCol{changeRow: 1, changeCol: -1}
	sync.ModifyArr[2] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[3] = &rowCol{changeRow: 1, changeCol: 0}
	return sync
}

//S:南
// 移动到南方向格子
func moveS() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 3, 3)
	sync.AppearArray[0] = &rowCol{changeRow: 2, changeCol: -1}
	sync.AppearArray[1] = &rowCol{changeRow: 2, changeCol: 0}
	sync.AppearArray[2] = &rowCol{changeRow: 2, changeCol: 1}

	sync.DisappearArray = make([]*rowCol, 3, 3)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.DisappearArray[1] = &rowCol{changeRow: -1, changeCol: 0}
	sync.DisappearArray[2] = &rowCol{changeRow: -1, changeCol: 1}

	sync.ModifyArr = make([]*rowCol, 6, 6)
	sync.ModifyArr[0] = &rowCol{changeRow: 1, changeCol: -1}
	sync.ModifyArr[1] = &rowCol{changeRow: 1, changeCol: 0}
	sync.ModifyArr[2] = &rowCol{changeRow: 1, changeCol: 1}
	sync.ModifyArr[3] = &rowCol{changeRow: 0, changeCol: -1}
	sync.ModifyArr[4] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[5] = &rowCol{changeRow: 0, changeCol: 1}
	return sync
}

//SE:东南
// 移动到东南方向格子
func moveSE() *SyncGirds {
	sync := &SyncGirds{}

	sync.AppearArray = make([]*rowCol, 5, 5)
	sync.AppearArray[0] = &rowCol{changeRow: 2, changeCol: 0}
	sync.AppearArray[1] = &rowCol{changeRow: 2, changeCol: 1}
	sync.AppearArray[2] = &rowCol{changeRow: 2, changeCol: 2}
	sync.AppearArray[3] = &rowCol{changeRow: 1, changeCol: 2}
	sync.AppearArray[4] = &rowCol{changeRow: 0, changeCol: 2}

	sync.DisappearArray = make([]*rowCol, 5, 5)
	sync.DisappearArray[0] = &rowCol{changeRow: -1, changeCol: -1}
	sync.DisappearArray[1] = &rowCol{changeRow: 0, changeCol: -1}
	sync.DisappearArray[2] = &rowCol{changeRow: 1, changeCol: -1}
	sync.DisappearArray[3] = &rowCol{changeRow: -1, changeCol: 0}
	sync.DisappearArray[4] = &rowCol{changeRow: -1, changeCol: 1}

	sync.ModifyArr = make([]*rowCol, 4, 4)
	sync.ModifyArr[0] = &rowCol{changeRow: 0, changeCol: 0}
	sync.ModifyArr[1] = &rowCol{changeRow: 0, changeCol: 1}
	sync.ModifyArr[2] = &rowCol{changeRow: 1, changeCol: 0}
	sync.ModifyArr[3] = &rowCol{changeRow: 1, changeCol: 1}
	return sync
}

func AddGridWatcher(gridComp *component.Grids, watchPos tools.Vec3f, RID int64, gateSession string) {
	gridKey := GetGridKey(watchPos)
	grid, ok := gridComp.Grids[gridKey]
	expect.True(ok, log.Fields{"gridKey": gridKey})
	if ok {
		grid.Watchers[RID] = true
	}
	w, ok := gridComp.Watchers[RID]
	if !ok {
		gridComp.Watchers[RID] = &component.Watcher{}
		w = gridComp.Watchers[RID]
	}
	w.GridKey = gridKey
	w.Session = gateSession
	w.WatchPos = watchPos
}
func DelGridWatcher(gridComp *component.Grids, RID int64) tools.Vec3f {
	watcher, ok := gridComp.Watchers[RID]
	if !ok {
		return tools.Invalid()
	}
	watchPos := watcher.WatchPos
	if g, ok := gridComp.Grids[watcher.GridKey]; ok {
		delete(g.Watchers, RID)
	}
	delete(gridComp.Watchers, RID)
	return watchPos
}

func AddEntityInGrid(w interfaces.IWorld, gridComp *component.Grids, gridKey int32, eid string) {
	ent := w.Runtime().GetEntity(eid)
	expect.True(ent != nil, log.Fields{"eid": eid})
	grid, ok := gridComp.Grids[gridKey]
	if ok {
		grid.Entities[ent.Id()] = ent
	}
}

func DelEntityInGrid(w interfaces.IWorld, gridComp *component.Grids, gridKey int32, eid string) {
	grid, ok := gridComp.Grids[gridKey]
	if ok {
		delete(grid.Entities, eid)
	}
}

func MoveEntityInGrid(w interfaces.IWorld, gridComp *component.Grids, oldKey, newKey int32, eid string) {
	if oldKey != newKey {
		return
	}
	DelEntityInGrid(w, gridComp, oldKey, eid)
	AddEntityInGrid(w, gridComp, newKey, eid)
}
