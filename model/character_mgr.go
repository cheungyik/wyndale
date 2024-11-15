package model

import "github.com/cornelk/hashmap"

// CharacterManager 角色管理器（绑定玩家与连接ID）
type CharacterManager struct {
	characterMapCID *hashmap.Map[int64, *Character]
}

func NewCharacterManager() *CharacterManager {
	return &CharacterManager{
		characterMapCID: hashmap.New[int64, *Character](),
	}
}

func (slf *CharacterManager) AddCharacter(character *Character) {
	slf.characterMapCID.Set(character.CID, character)
}

func (slf *CharacterManager) RemoveCharacter(character *Character) {
	slf.characterMapCID.Del(character.CID)
}

func (slf *CharacterManager) GetCharacterByCID(cid int64) *Character {
	v, ok := slf.characterMapCID.Get(cid)
	if ok {
		return v
	}
	return nil
}
