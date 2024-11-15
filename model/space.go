package model

import (
	"context"
	"github.com/cornelk/hashmap"
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
	"github.com/dobyte/due/v2/session"
	"wyndale/component/game/route"
	"wyndale/protocol/pb"
)

// Space 空间/地图/场景
type Space struct {
	ID          int64                           // ID
	Name        string                          // 名称
	Description string                          // 描述
	Music       string                          // 背景音乐
	characters  *hashmap.Map[int64, *Character] // 角色列表（实体ID与角色）
	Manager     *SpaceManager                   // 管理器
}

func NewSpace(id int64, name string) *Space {
	return &Space{
		ID:         id,
		Name:       name,
		characters: hashmap.New[int64, *Character](),
	}
}

func (slf *Space) app() *node.Proxy {
	return slf.Manager.app
}

func (slf *Space) CharacterJoin(character *Character) {
	log.Infof("角色 %v 加入场景 %v", character.GetEntityID(), slf.ID)
	// 设置角色所在场景
	character.SetSpace(slf)
	slf.characters.Set(character.entityID, character)

	notifyMsg := &pb.SpaceCharacterEnterNotify{
		SpaceId: slf.ID,
		EntityList: []*pb.Entity{
			character.ToProtoEntity(),
		},
	}
	// 向其它玩家广播有新角色加入场景
	var targets []int64
	slf.characters.Range(func(id int64, item *Character) bool {
		if id != character.entityID {
			targets = append(targets, item.CID)
		}
		return true
	})
	if len(targets) != 0 {
		if err := slf.app().Multicast(context.TODO(), &cluster.MulticastArgs{
			GID:     "gate",
			Kind:    session.Conn,
			Targets: targets,
			Message: &cluster.Message{
				Seq:   0,
				Route: route.SpaceCharactersEnterNotify,
				Data:  notifyMsg,
			},
		}); err != nil {
			log.Errorf("新角色进入场景时向其他玩家广播失败：%v", err)
		}
	}
	// 拉取场景内其它角色
	notifyMsg.EntityList = nil
	slf.characters.Range(func(id int64, item *Character) bool {
		if id != character.entityID {
			notifyMsg.EntityList = append(notifyMsg.EntityList, item.ToProtoEntity())
		}
		return true
	})
	if len(notifyMsg.EntityList) > 0 {
		if err := slf.app().Multicast(context.TODO(), &cluster.MulticastArgs{
			GID:     "gate",
			Kind:    session.Conn,
			Targets: []int64{character.CID},
			Message: &cluster.Message{
				Seq:   0,
				Route: route.SpaceCharactersEnterNotify,
				Data:  notifyMsg,
			},
		}); err != nil {
			log.Errorf("新角色进入场景时拉取已有其它角色失败：%v", err)
		}
	}
}

func (slf *Space) CharacterLeave(character *Character) {
	log.Infof("角色 %v 离开场景 %v", character.GetEntityID(), slf.ID)
	slf.characters.Del(character.entityID)
	notifyMsg := &pb.SpaceCharacterLeaveNotify{EntityId: character.GetEntityID()}
	var targets []int64
	slf.characters.Range(func(id int64, item *Character) bool {
		targets = append(targets, item.CID)
		return true
	})
	if len(targets) != 0 {
		if err := slf.app().Multicast(context.TODO(), &cluster.MulticastArgs{
			GID:     "gate",
			Kind:    session.Conn,
			Targets: targets,
			Message: &cluster.Message{
				Seq:   0,
				Route: route.SpaceCharactersLeaveNotify,
				Data:  notifyMsg,
			},
		}); err != nil {
			log.Errorf("角色离开场景时向其他玩家广播失败：%v", err)
		}
	}
}

func (slf *Space) UpdateCharacter(entitySync *pb.EntitySync) *pb.EntitySync {
	log.Infof("角色 %v 更新信息：%v", entitySync.GetEntity().GetId(), entitySync.GetEntity())
	var targets []int64
	slf.characters.Range(func(entityID int64, character *Character) bool {
		if entityID == entitySync.GetEntity().GetId() {
			character.Update4EntitySync(entitySync.GetEntity())
		} else {
			targets = append(targets, character.CID)
		}
		return true
	})
	character, exist := slf.characters.Get(entitySync.GetEntity().GetId())
	if !exist {
		return nil
	}
	updatedSync := &pb.EntitySync{
		Entity: character.ToProtoEntity(),
		State:  entitySync.GetState(),
	}
	if len(targets) != 0 {
		if err := slf.app().Multicast(context.TODO(), &cluster.MulticastArgs{
			GID:     "gate",
			Kind:    session.Conn,
			Targets: targets,
			Message: &cluster.Message{
				Seq:   0,
				Route: route.SpaceEntitySyncNotify,
				Data:  &pb.SpaceEntitySyncNotify{Sync: updatedSync},
			},
		}); err != nil {
			log.Errorf("角色更新信息时向其他玩家广播失败：%v", err)
		}
	}
	return updatedSync
}
