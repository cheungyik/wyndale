package logic

import (
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
	"math/rand/v2"
	"wyndale/component/game/route"
	"wyndale/model"
	"wyndale/protocol/pb"
)

type core struct {
	proxy        *node.Proxy
	entityMgr    *model.EntityManager
	spaceMgr     *model.SpaceManager
	characterMgr *model.CharacterManager // 玩家角色-UID
}

func NewCore(proxy *node.Proxy) *core {
	spaceMgr := model.NewSpaceManager(proxy)
	spaceMgr.AddSpace(model.NewSpace(1, "新手村"))
	return &core{
		proxy:        proxy,
		entityMgr:    model.NewEntityManager(),
		spaceMgr:     spaceMgr,
		characterMgr: model.NewCharacterManager(),
	}
}

func (slf *core) Init() {
	slf.proxy.Router().Group(func(group *node.RouterGroup) {
		group.AddRouteHandler(route.EnterGame, false, slf.enterGame)
		group.AddRouteHandler(route.SpaceEntitySync, false, slf.entitySync)
		group.AddRouteHandler(route.SpaceCharactersLeave, false, slf.leaveGame)
	})
	slf.proxy.AddEventHandler(cluster.Reconnect, slf.reconnect)
	slf.proxy.AddEventHandler(cluster.Disconnect, slf.disconnect)
}

func (slf *core) reconnect(ctx node.Context) {}
func (slf *core) disconnect(ctx node.Context) {
	character := slf.characterMgr.GetCharacterByCID(ctx.CID())
	if character == nil {
		return
	}
	space := character.GetSpace()
	if space == nil {
		return
	}
	space.CharacterLeave(character)
}

func (slf *core) enterGame(ctx node.Context) {
	req := &pb.GameEnterRequest{}
	res := &pb.GameEnterResponse{}
	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Success = false
		return
	}
	// 创建角色
	position, direction := &pb.Vector3{
		X: 500 + ((-5) + rand.Int64N(5-(-5))),
		Y: 0,
		Z: 500 + ((-5) + rand.Int64N(5-(-5))),
	}, &pb.Vector3{}
	position.X = position.X * 1000
	position.Y = position.Y * 1000
	position.Z = position.Z * 1000
	direction.X = direction.X * 1000
	direction.Y = direction.Y * 1000
	direction.Z = direction.Z * 1000
	character := model.NewCharacter(ctx.CID(), slf.entityMgr.NextEntityID(), position, direction)
	slf.characterMgr.AddCharacter(character)
	res.Entity = character.ToProtoEntity()
	res.Character = character.ToProtoCharacter()
	res.Success = true
	if err := ctx.Response(res); err != nil {
		log.Errorf("response message failed: %v", err)
		return
	}
	// 加入场景
	slf.spaceMgr.GetSpace(1).CharacterJoin(character)
}

func (slf *core) entitySync(ctx node.Context) {
	req := &pb.SpaceEntitySyncRequest{}
	res := &pb.SpaceEntitySyncResponse{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()
	character := slf.characterMgr.GetCharacterByCID(ctx.CID())
	if character == nil {
		return
	}
	space := character.GetSpace()
	if space == nil {
		return
	}
	res.Sync = space.UpdateCharacter(req.GetSync())
}

func (slf *core) leaveGame(ctx node.Context) {
	req := &pb.SpaceEntityLeaveRequest{}
	res := &pb.SpaceEntityLeaveResponse{}
	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		return
	}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()
	character := slf.characterMgr.GetCharacterByCID(ctx.CID())
	if character == nil {
		return
	}
	slf.characterMgr.RemoveCharacter(character)
	space := character.GetSpace()
	if space == nil {
		return
	}
	space.CharacterLeave(character)
}
