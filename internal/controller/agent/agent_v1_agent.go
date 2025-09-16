package agent

import (
	"agent/internal/service"
	"context"

	"agent/api/agent/v1"
)

func (c *ControllerV1) AgentStream(ctx context.Context, req *v1.AgentReq) (res *v1.AgentRes, err error) {
	return service.Agent().ReactAgentStream(ctx, req)
}
