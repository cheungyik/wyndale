package model

import "wyndale/protocol/pb"

type (
	// Position 位置
	Position = pb.Vector3
	// Direction 方向
	Direction = pb.Vector3
	// Entity 实体
	Entity struct {
		entityID  int64       // ID
		position  *pb.Vector3 // 位置
		direction *pb.Vector3 // 方向
		space     *Space      // 空间
	}
)

func NewEntity(id int64, position *Position, direction *Direction) *Entity {
	return &Entity{
		entityID:  id,
		position:  position,
		direction: direction,
	}
}

func (slf *Entity) GetEntityID() int64 {
	return slf.entityID
}

func (slf *Entity) SetPosition(position *Position) {
	slf.position = position
}

func (slf *Entity) GetPosition() *Position {
	return slf.position
}

func (slf *Entity) SetDirection(direction *Direction) {
	slf.direction = direction
}

func (slf *Entity) GetDirection() *Direction {
	return slf.direction
}

func (slf *Entity) ToProtoEntity() *pb.Entity {
	return &pb.Entity{
		Id:        slf.entityID,
		Position:  slf.position,
		Direction: slf.direction,
	}
}

func (slf *Entity) SetSpace(space *Space) {
	slf.space = space
}

func (slf *Entity) GetSpace() *Space {
	return slf.space
}

func (slf *Entity) Update4EntitySync(entity *pb.Entity) {
	slf.SetPosition(entity.GetPosition())
	slf.SetDirection(entity.GetDirection())
}
