package model

import "wyndale/protocol/pb"

// Character 角色
type Character struct {
	Actor
	CID int64
}

func NewCharacter(cid, entityID int64, position *Position, direction *Direction) *Character {
	return &Character{
		Actor: Actor{
			Entity: Entity{
				entityID:  entityID,
				position:  position,
				direction: direction,
			},
		},
		CID: cid,
	}
}

func (slf *Character) ToProtoCharacter() *pb.Character {
	return &pb.Character{}
}
