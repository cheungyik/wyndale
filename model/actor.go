package model

type Actor struct {
	Entity        // 继承 Entity
	Name   string // 名称
	Level  int    // 等级
}

func NewActor(id int64, position *Position, direction *Direction) *Actor {
	return &Actor{
		Entity: Entity{
			entityID:  id,
			position:  position,
			direction: direction,
		},
		Name: "",
	}
}
