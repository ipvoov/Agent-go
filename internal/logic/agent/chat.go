package agent

import (
	v1 "agent/api/agent/v1"
	"agent/internal/service"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	service.RegisterAgent(New())
}

type sAgent struct {
	historicalMessages map[string][]*schema.Message
}

func New() *sAgent {
	return &sAgent{
		historicalMessages: make(map[string][]*schema.Message),
	}
}

// ChainAgentStream 流式链式 Agent
func (s *sAgent) ChainAgentStream(ctx context.Context, in *v1.ChatStreamReq) {
	r := ghttp.RequestFromCtx(ctx)
	// 设置 SSE 响应头
	r.Response.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	chatModel := NewChatModel(ctx)

	template, example := Template(ctx, &v1.ChatStreamReq{
		Query:     in.Query,
		SessionID: in.SessionID,
	})
	sessionMessages := s.GetSessionBySessionID(ctx, in.SessionID)
	variables := map[string]any{
		"role":        "expert in the field of relationships with many years of experience",
		"example":     example,
		"task":        in.Query,
		"history_key": sessionMessages,
	}
	messages, _ := template.Format(ctx, variables)

	reader, err := chatModel.Stream(ctx, messages)
	if err != nil {
		SndErr(r, err)
	}
	defer reader.Close()

	r.Response.WriteHeader(200)
	r.Response.Flush()

	var fullContent strings.Builder
	var chunk *schema.Message
	for {
		chunk, err = reader.Recv()
		if err != nil {
			r.Response.Write([]byte("data: {\"content\":\"\",\"done\":true}\n\n"))
			r.Response.Flush()
			break
		}
		fullContent.Write([]byte(chunk.Content))
		resp := v1.ChatStreamRes{
			Content: chunk.Content,
			Done:    false,
		}
		data, _ := gjson.Marshal(resp)
		r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
		r.Response.Flush()
	}

	finalContent := fullContent.String()
	s.historicalMessages[in.SessionID] = append(s.historicalMessages[in.SessionID],
		schema.UserMessage(in.Query),
		schema.AssistantMessage(finalContent, nil))

	if err != nil && err != io.EOF {
		SndErr(r, err)
	}
	return
}

func SndErr(r *ghttp.Request, err error) {
	// 发送错误信息
	errorData := g.Map{
		"error": err.Error(),
		"done":  true,
	}
	jsonData, _ := gjson.Marshal(errorData)
	r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonData)))
	r.Response.Flush()
}

// GetSessionBySessionID 获取会话
func (s *sAgent) GetSessionBySessionID(ctx context.Context, sessionId string) []*schema.Message {
	return s.historicalMessages[sessionId]
}

// IsRelevant 检查检索结果是否与查询相关
func (s *sAgent) IsRelevant(query string, results []*schema.Document) bool {
	// 如果查询太短（少于3个字符），要求更高的相关性
	if len(strings.TrimSpace(query)) < 10 {
		return false
	}

	// 检查是否包含一些关键词匹配
	queryLower := strings.ToLower(query)
	for _, doc := range results {
		contentLower := strings.ToLower(doc.Content)
		// 如果文档内容包含查询的关键词，认为相关
		if strings.Contains(contentLower, queryLower) {
			return true
		}
	}

	return false
}
