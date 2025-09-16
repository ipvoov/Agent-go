package agent

import (
	"context"

	v1 "agent/api/agent/v1"
	"agent/internal/service"
)

func (c *ControllerV1) ChatStream(ctx context.Context, in *v1.ChatStreamReq) (out *v1.ChatStreamRes, err error) {
	service.Agent().ChainAgentStream(ctx, in)
	return out, nil
}
