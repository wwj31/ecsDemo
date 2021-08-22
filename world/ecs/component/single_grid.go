package component

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/tools"
)

type (
	Grids struct {
		ecs.ComponentBase
		Watchers map[int64]*Watcher //RId =>
		Grids    map[int32]*Grid    //  x/size<<16+y/size
	}
	Grid struct {
		Entities map[string]*ecs.Entity
		//NeedSync bool
		Watchers map[int64]bool // //RId =>
	}
	Watcher struct {
		GridKey  int32
		WatchPos tools.Vec3f
		Session  string
	}
)

func NewGrids() *Grids {
	v := &Grids{
		Watchers: make(map[int64]*Watcher),
		Grids:    make(map[int32]*Grid),
	}
	return v
}

func (s *Grids) Type() ecs.ComponentType {
	return GRID_COMP
}
