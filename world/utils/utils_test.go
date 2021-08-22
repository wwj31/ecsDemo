package utils

import (
	"fmt"
	"github.com/wwj31/dogactor/tools"
	"math/rand"
	"ecsDemo/world/constant"
	"sort"
	"testing"
)

func TestUtils(t *testing.T) {
	//i, x, y := GridIndex(2, 2, 1000, 1001, tools.Vec3f{X: 499, Y: 499})
	//println(i, x, y)
	fmt.Println(ActionWithMove(EnterGrid, 1, 1))

}

func Test_MapView(t *testing.T) {
	rand.Seed(tools.Now().Unix())
	for i := 0; i < 1000000; i++ {
		fromX := float64(rand.Intn(constant.MAX_WORLD_WIDTH))
		fromY := float64(rand.Intn(constant.MAX_WORLD_HEIGHT))
		var toX, toY float64
		if rand.Intn(10000) > 5000 {
			toX = fromX + constant.GRID_SIZE
		} else {
			toX = fromX - constant.GRID_SIZE
		}
		if rand.Intn(10000) > 5000 {
			toY = fromY + constant.GRID_SIZE
		} else {
			toY = fromY - constant.GRID_SIZE
		}

		appear1, disAppear1, modify1 := SyncGrids(fromX, fromY, toX, toY)
		sort.Slice(appear1, func(i, j int) bool {
			return appear1[i] < appear1[j]
		})
		sort.Slice(disAppear1, func(i, j int) bool {
			return disAppear1[i] < disAppear1[j]
		})
		sort.Slice(modify1, func(i, j int) bool {
			return modify1[i] < modify1[j]
		})

		appear2, disAppear2, modify2 := SyncGrids2(fromX, fromY, toX, toY)

		sort.Slice(appear2, func(i, j int) bool {
			return appear2[i] < appear2[j]
		})
		sort.Slice(disAppear2, func(i, j int) bool {
			return disAppear2[i] < disAppear2[j]
		})
		sort.Slice(modify2, func(i, j int) bool {
			return modify2[i] < modify2[j]
		})
		var flag bool
		if !isEquip(appear1, appear2) {
			flag = true
			fmt.Println("appear1:", appear1)
			for _, v := range appear1 {
				fmt.Println(tools.Int32Split(v))
			}
			fmt.Println("appear2:", appear2)
			for _, v := range appear2 {
				fmt.Println(tools.Int32Split(v))
			}
		}
		if !isEquip(disAppear1, disAppear2) {
			fmt.Println("disAppear1:", disAppear1)
			for _, v := range disAppear1 {
				fmt.Println(tools.Int32Split(v))
			}
			fmt.Println("disAppear2:", disAppear2)
			for _, v := range disAppear2 {
				fmt.Println(tools.Int32Split(v))
			}
			flag = true
		}
		if !isEquip(modify1, modify2) {
			fmt.Println("modify1:", modify1)
			for _, v := range modify1 {
				fmt.Println(tools.Int32Split(v))
			}
			fmt.Println("modify2:", modify2)
			for _, v := range modify2 {
				fmt.Println(tools.Int32Split(v))
			}
			flag = true
		}
		if flag {
			fromRow, fromCol := GetGridRowCol(fromX, fromY)
			toRow, toCol := GetGridRowCol(toX, toY)
			fmt.Println("------", fromX, fromY, toX, toY)
			fmt.Println(fromRow, fromCol, toRow, toCol)
			SyncGrids(fromX, fromY, toX, toY)
			SyncGrids2(fromX, fromY, toX, toY)
			return
		}
	}
}

func isEquip(arr1 []int32, arr2 []int32) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for index, _ := range arr1 {
		if arr1[index] != arr2[index] {
			return false
		}
	}
	return true
}

//BenchmarkSyncGrids
//BenchmarkSyncGrids-8    	 2546557	       470 ns/op
//BenchmarkSyncGrids2
//BenchmarkSyncGrids2-8   	  433249	      2758 ns/op
//PASS

func BenchmarkSyncGrids(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fromX := float64(rand.Intn(1200))
		fromY := float64(rand.Intn(1200))
		var toX, toY float64
		if rand.Intn(10000) > 5000 {
			toX = fromX + constant.GRID_SIZE
		} else {
			toX = fromX - constant.GRID_SIZE
		}
		if rand.Intn(10000) > 5000 {
			toY = fromY + constant.GRID_SIZE
		} else {
			toY = fromY - constant.GRID_SIZE
		}
		SyncGrids(fromX, fromY, toX, toY)
	}
}
func BenchmarkSyncGrids2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fromX := float64(rand.Intn(constant.MAX_WORLD_WIDTH))
		fromY := float64(rand.Intn(constant.MAX_WORLD_HEIGHT))
		var toX, toY float64
		if rand.Intn(10000) > 5000 {
			toX = fromX + constant.GRID_SIZE
		} else {
			toX = fromX - constant.GRID_SIZE
		}
		if rand.Intn(10000) > 5000 {
			toY = fromY + constant.GRID_SIZE
		} else {
			toY = fromY - constant.GRID_SIZE
		}
		SyncGrids2(fromX, fromY, toX, toY)
	}
}
