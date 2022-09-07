package main

import (
	_ "github.com/morgeq/iotfast/plugin/service"
	plugin "github.com/morgeq/iotfast/plugin/server"
	
)

func main() {
	plu := plugin.NewServer("unix", "E:\\src\\github.com/morgeq/iotfast\\plugin\\examples\\client")
	defer plu.Stop()
	plu.Start()

	// cfg := plugin.PluginCfg{
	// 	Proto:      "tcp",
	// 	Unixdir:    "",
	// 	ServerAddr: "127.0.0.1:1024",
	// 	Name:       "client",
	// 	Id:         string(libUtils.GetRandomUUID()),
	// 	Params:     "",
	// }
	// p := plu.AddPlugin("G:\\OCM\\private\\github.com/morgeq/iotfast\\plugin\\examples\\client\\client.exe", cfg)
	// plu.PluginStart(p)
}
