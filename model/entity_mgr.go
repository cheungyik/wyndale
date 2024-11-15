package model

import (
	"github.com/cornelk/hashmap"
	"sync"
	"wyndale/protocol/pb"
)

type EntityManager struct {
	entityIndex       int64
	entityIndexLocker sync.Mutex
	entities          *hashmap.Map[int64, *Entity]
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entityIndex:       0,
		entityIndexLocker: sync.Mutex{},
		entities:          hashmap.New[int64, *Entity](),
	}
}

func (slf *EntityManager) NextEntityID() int64 {
	slf.entityIndexLocker.Lock()
	defer slf.entityIndexLocker.Unlock()
	slf.entityIndex++
	return slf.entityIndex
}

func (slf *EntityManager) CreateEntity() *Entity {
	entity := NewEntity(slf.NextEntityID(), &pb.Vector3{}, &pb.Vector3{})
	slf.entities.Set(entity.GetEntityID(), entity)
	return entity
}
