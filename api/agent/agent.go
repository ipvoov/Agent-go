// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package agent

import (
	"context"

	"agent/api/agent/v1"
)

type IAgentV1 interface {
	ChatStream(ctx context.Context, req *v1.ChatStreamReq) (res *v1.ChatStreamRes, err error)
	AgentStream(ctx context.Context, req *v1.AgentReq) (res *v1.AgentRes, err error)
}
