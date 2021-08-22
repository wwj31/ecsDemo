package graphUI

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/world/constant"
	"fmt"
	"github.com/gonutz/prototype/draw"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"math"
	"os"
	"time"
)

var uiActor *GraphUI
var exit chan os.Signal

func Init(exitCh chan os.Signal) {
	exit = exitCh
	draw.RunWindow("map", constant.MAX_WORLD_WIDTH, constant.MAX_WORLD_HEIGHT, update)
}

func update(window draw.Window) {
	if uiActor == nil {
		return
	}

	select {
	case <-exit:
		window.Close()
		return
	default:
		//continue
	}

	window.FillRect(0, 0, constant.MAX_WORLD_WIDTH, constant.MAX_WORLD_HEIGHT, draw.LightGray)

	uiActor.l.Lock()
	defer uiActor.l.Unlock()

	for _, p := range uiActor.Path {
		window.FillEllipse(int(p.X), int(p.Y), 5, 5, draw.DarkRed)
	}

	// 绘hp，圆
	for k, enti := range uiActor.entities {
		dv := area(k)
		val := 20
		posoffset := val / 2
		x := int(enti.X - float64(posoffset))
		y := int(enti.Y - float64(posoffset))
		window.FillEllipse(x, y, val, val, colorAreaNum[dv])
		if !enti.Duplicate {
			window.DrawText(fmt.Sprintf("%v", enti.HP), x+2, y+2, draw.White)
		}
	}

	// 绘区域
	for _, t := range uiActor.tiles {
		window.DrawRect(
			int(t.ExclusiveBound.Pos.X), int(t.ExclusiveBound.Pos.Y),
			int(t.ExclusiveBound.Width), int(t.ExclusiveBound.Height),
			draw.RGB(0, 200, 200),
		)
		window.DrawRect(
			int(t.ActualBound.Pos.X), int(t.ActualBound.Pos.Y),
			int(t.ActualBound.Width), int(t.ActualBound.Height),
			draw.White,
			//draw.RGB(60, 60, 60),
		)
		window.DrawText(fmt.Sprintf("%v", t.Count), int(140+t.ActualBound.Pos.X), int(140+t.ActualBound.Pos.Y), draw.RGB(0, 0, 0))
	}

	mouseX, mouseY := window.MousePosition()
	mouseInCircle := math.Hypot(float64(mouseX), float64(mouseY)) < 40
	color := draw.DarkRed
	if mouseInCircle {
		color = draw.Red
	}
	window.FillEllipse(0, 0, 40, 40, color)
	window.DrawEllipse(0, 0, 40, 40, draw.White)
	if mouseInCircle {
		window.DrawScaledText("random!", 40, 25, 1.6, draw.DarkBlue)
	}

	// check all mouse clicks that happened during this frame
	for _, click := range window.Clicks() {
		aNum := areaNum(click.X, click.Y)
		switch click.Button {
		case draw.LeftButton:
			mouseInCircle = math.Hypot(float64(click.X), float64(click.Y)) < 40
			if mouseInCircle {
				uiActor.random()
			}

			entity := uiActor.entities[foramtEid(uiActor.lastEID, int64(aNum))]
			if entity == nil {
				return
			}
			entity.Path = uiActor.Path
			if len(entity.Path) > 0 {
				aNum = areaNum(int(entity.X), int(entity.Y))
				log.KVs(log.Fields{"aNum": aNum}).Debug("move")
				uiActor.Send(common.WorldName(int32(aNum)), &inner.U2GMoveEntity{EID: uiActor.lastEID, Path: entity.Path, Speed: 35})

				uiActor.Path = entity.Path
				entity.Path = make([]*inner.Vec3F, 0)
				uiActor.clear = 1
			}

		case draw.RightButton:
			if uiActor.clear == 1 {
				uiActor.Path = make([]*inner.Vec3F, 0)
				uiActor.clear = 0
			}
			uiActor.Path = append(uiActor.Path, &inner.Vec3F{X: float64(click.X), Y: float64(click.Y)})
		case draw.MiddleButton:
			for i := 0; i < 1; i++ {
				time.Sleep(10 * time.Millisecond)
				//randpos := &inner.Vec3F{X: float64(rand.Intn(850) + 1), Y: float64(rand.Intn(850) + 1)}
				uiActor.lastEID = tools.UUID()
				randpos := &inner.Vec3F{X: float64(click.X), Y: float64(click.Y)}
				tempa := areaNum(int(randpos.X), int(randpos.Y))
				uiActor.Send(common.WorldName(int32(tempa)),
					&inner.U2GAddEntity{
						AreaNum: int32(aNum),
						EID:     uiActor.lastEID,
						RealPos: randpos,
					})
			}
		}
	}
}

var colorAreaNum = map[int]draw.Color{
	0: draw.DarkYellow,
	1: draw.Purple,
	2: draw.Green,
	3: draw.Red,
	4: draw.Brown,
	5: draw.Blue,
	6: draw.DarkPurple,
	7: draw.Cyan,
	8: draw.DarkGreen,
}
