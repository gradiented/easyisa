package main

import (
	cadence "github.com/gradiented/easyisa/cmd/cadence"
	"github.com/gradiented/easyisa/pkg/server"
)

func main() {
	cadence.Start()
	server.Start()
}
