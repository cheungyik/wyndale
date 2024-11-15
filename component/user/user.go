package user

import (
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/registry/consul/v2"
	"github.com/dobyte/due/v2/cluster/node"
	"wyndale/component/user/logic"
)

func NewComponent() *node.Node {
	component := node.NewNode(
		node.WithName("user"),
		node.WithLocator(redis.NewLocator()),
		node.WithRegistry(consul.NewRegistry()),
	)
	logic.NewCore(component.Proxy()).Init()
	return component
}
