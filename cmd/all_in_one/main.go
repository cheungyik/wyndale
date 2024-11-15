package main

import (
	"github.com/dobyte/due/v2"
	"wyndale/component/game"
	"wyndale/component/user"
)

func main() {
	container := due.NewContainer()
	container.Add(user.NewComponent(), game.NewComponent())
	container.Serve()
}
