package model

import (
	"github.com/cornelk/hashmap"
	"github.com/dobyte/due/v2/cluster/node"
)

type SpaceManager struct {
	app    *node.Proxy // 应用
	spaces *hashmap.Map[int64, *Space]
}

func NewSpaceManager(app *node.Proxy) *SpaceManager {
	return &SpaceManager{
		app:    app,
		spaces: hashmap.New[int64, *Space](),
	}
}

func (slf *SpaceManager) GetSpace(id int64) *Space {
	v, ok := slf.spaces.Get(id)
	if ok {
		return v
	}
	return nil
}

func (slf *SpaceManager) AddSpace(space *Space) {
	space.Manager = slf
	slf.spaces.Set(space.ID, space)
}
