package utils

//对比用
func SyncGrids2(fromX float64, fromY float64, toX float64, toY float64) (appearArray []int32, disappearArr []int32, modifyArr []int32) {
	fromRow,fromCol :=GetGridRowCol(fromX,fromY)
	toRow, toCol := GetGridRowCol(toX, toY)
	switch {
	case fromX == -1 && fromY == -1: //进入
		curArea := GetSurroundGridsByGid(fromRow,fromCol)
		for curGridId, _ := range curArea {
			appearArray = append(appearArray, curGridId) //修改
		}
	case toX == -1 && toY == -1: //离开
		curArea := GetSurroundGridsByGid(fromRow,fromCol)
		for curGridId, _ := range curArea {
			disappearArr = append(disappearArr, curGridId) //修改
		}
	case (fromRow == toRow) && (fromCol==toCol):
		curArea := GetSurroundGridsByGid(fromRow,fromCol)
		for curGridId, _ := range curArea {
			modifyArr = append(modifyArr, curGridId) //修改
		}
	default:
		oldArea := GetSurroundGridsByGid(fromRow,fromCol)
		curArea := GetSurroundGridsByGid(toRow,toCol)
		for oldGridId, _ := range oldArea {
			if _, exist := curArea[oldGridId]; !exist {
				disappearArr = append(disappearArr, oldGridId) //消失
			}
		}
		for curGridId, _ := range curArea {
			if _, exist := oldArea[curGridId]; !exist {
				appearArray = append(appearArray, curGridId) //出现
			} else {
				modifyArr = append(modifyArr, curGridId) //修改
			}
		}
	}
	return
}
var _surroundGrids =make([]*rowCol,9,9)
func init(){
	_surroundGrids = make([]*rowCol,9,9)
	_surroundGrids[0]=&rowCol{changeRow:-1,changeCol:-1} //NW
	_surroundGrids[1]=&rowCol{changeRow:-1,changeCol:0}  //N
	_surroundGrids[2]=&rowCol{changeRow:-1,changeCol:1} ////NE
	_surroundGrids[3]=&rowCol{changeRow:0,changeCol:-1} //W
	_surroundGrids[4]=&rowCol{changeRow:0,changeCol:0} //self
	_surroundGrids[5]=&rowCol{changeRow:0,changeCol:1} //E
	_surroundGrids[6]=&rowCol{changeRow:1,changeCol:-1} //SW
	_surroundGrids[7]=&rowCol{changeRow:1,changeCol:0} //S
	_surroundGrids[8]=&rowCol{changeRow:1,changeCol:1} //SE
}


//根据格子的gID得到当前周边的九宫格信息
func  GetSurroundGridsByGid(row int,col int) (grids map[int32]bool) {
	//将当前gid添加到九宫格中
	grids = make(map[int32]bool)
	for _,v:= range _surroundGrids{
		gridId:=GetGridByRowCol(v.changeRow+row,v.changeCol+col)
		if gridId < 0{
			continue
		}
		grids[gridId] =true
	}
	return
}