package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ChatStreamReq struct {
	g.Meta    `path:"/chatSteam" method:"get" summary:"You first agent api"`
	Query     string `json:"query" p:"query" v:"required"`
	SessionID string `json:"session_id" p:"session_id" v:"required"`
}

type ChatStreamRes struct {
	Content  string `json:"content"`
	Thinking bool   `json:"thinking"`
	Done     bool   `json:"done"`
}

type AgentReq struct {
	g.Meta    `path:"/agentStream"  method:"get" summary:"You first agent api"`
	Query     string `json:"query" p:"query" v:"required"`
	SessionID string `json:"session_id" p:"session_id" v:"required"`
}
type AgentRes struct {
	Content      string `json:"content"`
	TotalTokens  int    `json:"total_tokens"`
	MessageCount int    `json:"message_count"`
}
