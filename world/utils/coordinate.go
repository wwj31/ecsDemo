package utils

import (
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"ecsDemo/world/constant"
)

// 传入位置，得到所在区域
func GridIndex(pos tools.Vec3f) (i, x, y int32) {
	return gridIndex(constant.SERVER_SPLIT_AREA, constant.SERVER_SPLIT_AREA, constant.MAX_WORLD_WIDTH, constant.MAX_WORLD_HEIGHT, pos)
}

/*
	传入一个坐标点，得到在宽*高网格内的索引				0 1 2
	wg,hg 几×几的网格					 索引排列  = 3 4 5
	wt,ht 宽高总长度									6 7 8
	pos 网格内的点
*/
func gridIndex(wg, hg int32, wt, ht float64, pos tools.Vec3f) (index int32, x, y int32) {
	if pos.X > wt || pos.Y > ht || pos.X < 0 || pos.Y < 0 {
		log.KVs(log.Fields{"pos": pos, "wt": wt, "ht": ht}).ErrorStack(3, "out of bound")
		return -1, -1, -1
	}
	w, h := wt/float64(wg), ht/float64(hg) // 单位格子宽、高
	x = int32(pos.X / w)
	y = int32(pos.Y / h)

	return y*wg + x, x, y
}

func GetGridKey(p tools.Vec3f) int32 {
	return tools.Int32Merge(int16(p.X/constant.GRID_SIZE), int16(p.Y/constant.GRID_SIZE))
}
func DivGridKey(gridKey int32) tools.Vec3f {
	h, l := tools.Int32Split(gridKey)
	return tools.Vec3f{X: float64(h), Y: float64(l)}
}

func GetGridRowCol(x float64, y float64) (row int, col int) {
	return int(x / constant.GRID_SIZE), int(y / constant.GRID_SIZE)
}

func GetGridByRowCol(row int, col int) int32 {
	return tools.Int32Merge(int16(row), int16(col))
}
