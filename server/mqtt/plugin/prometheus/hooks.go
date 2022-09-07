package prometheus

import (
	"github.com/morgeq/iotfast/server/mqtt/server"
)

func (p *Prometheus) HookWrapper() server.HookWrapper {
	return server.HookWrapper{}
}
