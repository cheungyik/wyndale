package main

import (
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/network/tcp/v2"
	"github.com/dobyte/due/registry/consul/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/gate"
)

func main() {
	container := due.NewContainer()
	component := gate.NewGate(
		gate.WithServer(tcp.NewServer()),
		gate.WithLocator(redis.NewLocator()),
		gate.WithRegistry(consul.NewRegistry()),
	)
	container.Add(component)
	container.Serve()
}
