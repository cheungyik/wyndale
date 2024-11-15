package main

import (
	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/client"
	"wyndale/cluster/client/app"
)

func main() {
	container := due.NewContainer()
	component := client.NewClient(
		client.WithClient(ws.NewClient()),
	)
	app.Init(component.Proxy())
	container.Add(component)
	container.Serve()
}
