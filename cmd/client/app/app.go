package app

import (
	"github.com/dobyte/due/v2/cluster/client"
	"wyndale/cmd/client/app/logic"
)

func Init(proxy *client.Proxy) {
	logic.NewCode(proxy).Init()
}
