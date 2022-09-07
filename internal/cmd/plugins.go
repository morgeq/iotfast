package cmd

import (
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/morgeq/iotfast/plugin/service"
	_ "github.com/morgeq/iotfast/server/mqtt/plugin/admin"
	_ "github.com/morgeq/iotfast/server/mqtt/plugin/auth"
	_ "github.com/morgeq/iotfast/server/mqtt/plugin/device"
	_ "github.com/morgeq/iotfast/server/mqtt/plugin/prometheus"
)
