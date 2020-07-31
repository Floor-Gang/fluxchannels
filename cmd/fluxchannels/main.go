package main

import (
	"github.com/Floor-Gang/fluxchannels/internal"
	util "github.com/Floor-Gang/utilpkg"
)

func main() {
	internal.Start()
	util.KeepAlive()
}
