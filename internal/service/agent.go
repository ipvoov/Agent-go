// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "agent/api/agent/v1"
	"context"

	"github.com/cloudwego/eino/schema"
)

type (
	IAgent interface {
		ReactAgentStream(ctx context.Context, in *v1.AgentReq) (out *v1.AgentRes, err error)
		// ChainAgentStream 流式链式 Agent
		ChainAgentStream(ctx context.Context, in *v1.ChatStreamReq)
		// GetSessionBySessionID 获取会话
		GetSessionBySessionID(ctx context.Context, sessionId string) []*schema.Message
		// IsRelevant 检查检索结果是否与查询相关
		IsRelevant(query string, results []*schema.Document) bool
	}
)

var (
	localAgent IAgent
)

func Agent() IAgent {
	if localAgent == nil {
		panic("implement not found for interface IAgent, forgot register?")
	}
	return localAgent
}

func RegisterAgent(i IAgent) {
	localAgent = i
}
