package model

// Monster 怪物
type Monster struct {
	Actor
}

func NewMonster(id int64, position *Position, direction *Direction) *Monster {
	return &Monster{
		Actor: Actor{
			Entity: Entity{
				entityID:  id,
				position:  position,
				direction: direction,
			},
		},
	}
}
