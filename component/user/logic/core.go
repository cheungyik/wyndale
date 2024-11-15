package logic

import (
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
	"wyndale/codes"
	"wyndale/component/user/route"
	"wyndale/middleware"
	"wyndale/model"
	"wyndale/protocol/pb"
)

type core struct {
	proxy            *node.Proxy
	usersMapID       map[int64]*model.User
	usersMapUsername map[string]*model.User
	uidIndex         int64
}

func (slf *core) genUID() int64 {
	slf.uidIndex++
	return slf.uidIndex
}

func NewCore(proxy *node.Proxy) *core {
	return &core{
		proxy:            proxy,
		usersMapID:       make(map[int64]*model.User),
		usersMapUsername: make(map[string]*model.User),
	}
}

func (slf *core) Init() {
	slf.proxy.Router().Group(func(group *node.RouterGroup) {
		group.AddRouteHandler(route.Register, false, slf.register)
		group.AddRouteHandler(route.Login, false, slf.login)
		group.Middleware(middleware.Auth)
	})
	slf.proxy.AddEventHandler(cluster.Reconnect, slf.reconnect)
	slf.proxy.AddEventHandler(cluster.Disconnect, slf.disconnect)
}

func (slf *core) reconnect(ctx node.Context)  {}
func (slf *core) disconnect(ctx node.Context) {}

func (slf *core) register(ctx node.Context) {
	req := &pb.UserRegisterRequest{}
	res := &pb.UserRegisterResponse{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()
	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = int32(codes.InternalError.Code())
		return
	}
	if _, ok := slf.usersMapUsername[req.Username]; ok {
		res.Code = int32(codes.AccountExists.Code())
		return
	}
	user := &model.User{
		ID:       slf.genUID(),
		Username: req.Username,
		Password: req.Password,
	}
	slf.usersMapID[user.ID] = user
	slf.usersMapUsername[user.Username] = user
	res.Code = int32(codes.OK.Code())
}

func (slf *core) login(ctx node.Context) {
	req := &pb.UserLoginRequest{}
	res := &pb.UserLoginResponse{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()
	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = int32(codes.InternalError.Code())
		return
	}
	user, ok := slf.usersMapUsername[req.Username]
	if !ok || user.Password != req.Password {
		res.Code = int32(codes.AccountNotExists.Code())
		return
	}
	if err := ctx.BindGate(user.ID); err != nil {
		log.Errorf("bind gate failed: %v", err)
		res.Code = int32(codes.InternalError.Code())
		return
	}
	res.Code = int32(codes.OK.Code())
}
