package main

import (
	_ "github.com/morgeq/iotfast/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/morgeq/iotfast/internal/cmd"
)

func main() {
	// cmd.Main.Run(gctx.New())

	err := cmd.Main.AddCommand(cmd.Plugin)
	if err != nil {
		panic(err)
	}
	cmd.Main.Run(gctx.New())

}
