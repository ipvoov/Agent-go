package agent

import (
	v1 "agent/api/agent/v1"
	"agent/internal/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (s *sAgent) ReactAgentStream(ctx context.Context, in *v1.AgentReq) (out *v1.AgentRes, err error) {
	chatModel := NewChatModel(ctx)

	toolCallChecker := func(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
		defer sr.Close()
		for {
			msg, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					// finish
					break
				}

				return false, err
			}

			if len(msg.ToolCalls) > 0 {
				return true, nil
			}
		}
		return false, nil
	}

	raAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				tools.NewPDFGenerationTool(),
				tools.NewWebSearchTool(),
				tools.NewResourceDownloadTool(),
				tools.NewPhotoSearchTool(),
				//tools.NewFileOperationTool(),
				//tools.NewTerminalOperationTool(),
				//GetGdMapMCPTool(ctx),
			},
			ExecuteSequentially: false,
		},
		StreamToolCallChecker: toolCallChecker,
	})
	if err != nil {
		return
	}

	template := AgentTemplate(ctx, &v1.ChatStreamReq{
		Query:     in.Query,
		SessionID: in.SessionID,
	})

	r := ghttp.RequestFromCtx(ctx)
	// 设置 SSE 响应头
	r.Response.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	r.Response.WriteHeader(200)
	r.Response.Flush()

	resp, err := raAgent.Stream(ctx, template, agent.WithComposeOptions(compose.WithCallbacks(&LoggerCallback{})))
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	data := ctx.Value("data")
	if data != nil {
		for _, v := range data.(string) {
			resps := v1.ChatStreamRes{
				Content:  string(v),
				Thinking: false,
				Done:     true,
			}
			d, _ := gjson.Marshal(resps)
			r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
			r.Response.Flush()
			time.Sleep(100 * time.Microsecond)
		}
	}
	respss := v1.ChatStreamRes{
		Content:  "",
		Thinking: false,
		Done:     true,
	}
	d, _ := gjson.Marshal(respss)
	r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
	r.Response.Flush()
	out = &v1.AgentRes{}
	return
}

type LoggerCallback struct {
	callbacks.HandlerBuilder // 可以用 callbacks.HandlerBuilder 来辅助实现 callback
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	//fmt.Println("==================")
	//inputStr, _ := gjson.MarshalIndent(input, "", "  ")
	//fmt.Printf("[OnStart] %s\n", string(inputStr))
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	//fmt.Println("=========[OnEnd]=========")
	//outputStr, _ := gjson.MarshalIndent(output, "", "  ")

	//fmt.Println(string(outputStr))
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========[OnError]=========")
	fmt.Println(err, "\ninfo.Name:", info.Name)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	r := ghttp.RequestFromCtx(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	var tol string
	go func(tol *string) {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[OnEndStream] panic err:", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer
		fmt.Println("=========[OnEndStream]=========")
		tol2 := ""
		for {
			if output == nil {
				return
			}
			frame, err := output.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}
			s, _ := json.Marshal(frame)
			if info.Name == react.ToolsNodeName {
				str := fmt.Sprintf(string(s))
				var toolResp []ToolResponse
				err = gjson.DecodeTo(str, &toolResp)
				if err != nil {
					fmt.Printf("internal error: %s\n", err)
					return
				}
				var data2 []SearchResult
				if len(toolResp) == 0 {
					fmt.Println("warning: toolResp is empty")
					return
				}
				err = gjson.DecodeTo(toolResp[0].Content, &data2)
				if err != nil {
					return
				}
				for _, data3 := range data2 {
					resp := v1.ChatStreamRes{
						Content: fmt.Sprintf("Date: %s\nDisplayedLink: %s\nLink: %s\nPosition: %d\nSnippet: %s\nSnippetHighlightedWords: %s\nThumbnail: %s\nTitle: %s\n\n",
							data3.Date, data3.DisplayedLink, data3.Link, data3.Position, data3.Snippet, data3.SnippetHighlightedWords, data3.Thumbnail, data3.Title),
						Thinking: true,
						Done:     false,
					}
					// 序列化响应对象
					d, _ := gjson.Marshal(resp)
					// 使用 SSE 格式发送数据
					r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
					r.Response.Flush() // 立即发送数据
					fmt.Printf(resp.Content)
					time.Sleep(100 * time.Microsecond)
				}
				return
			} else if info.Name == react.ModelNodeName {
				var data CallbackOutput
				err = gjson.DecodeTo(s, &data)
				if err != nil {
					fmt.Printf("%s", string(s))
					str := "\n" + string(s) + "\n\n\n"
					for _, char := range str {
						resp := v1.ChatStreamRes{
							Content:  string(char),
							Thinking: true,
							Done:     false,
						}
						d, _ := gjson.Marshal(resp)
						r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
						r.Response.Flush()
					}
					return
				}
				resp := v1.ChatStreamRes{
					Content:  data.Message.Content,
					Thinking: true,
					Done:     false,
				}
				d, _ := gjson.Marshal(resp)
				r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
				r.Response.Flush()
				*tol += data.Message.Content
				tol2 += data.Message.Content
			} else if info.Name == react.GraphName {
				var message schema.Message
				err = gjson.DecodeTo(s, &message)
				if err != nil {
					fmt.Printf("%s", string(s))
					str := "\n" + string(s) + "\n"
					for _, char := range str {
						resp := v1.ChatStreamRes{
							Content:  string(char),
							Thinking: false,
							Done:     false,
						}
						d, _ := gjson.Marshal(resp)
						r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
						r.Response.Flush()
					}
					return
				}
				resp := v1.ChatStreamRes{
					Content:  message.Content,
					Thinking: false,
					Done:     false,
				}
				d, _ := gjson.Marshal(resp)
				r.Response.Write([]byte(fmt.Sprintf("data: %s\n\n", d)))
				r.Response.Flush()
				*tol += message.Content
				tol2 += message.Content
			}
		}
		fmt.Printf("%s:%s\n", info.Name, tol2)
	}(&tol)
	wg.Wait()
	ctx = context.WithValue(ctx, "data", tol)
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}

type ToolResponse struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id"`
	ToolName   string `json:"tool_name"`
}

type SearchResult struct {
	Date                    string   `json:"date"`
	DisplayedLink           string   `json:"displayed_link"`
	Link                    string   `json:"link"`
	Position                int      `json:"position"`
	Snippet                 string   `json:"snippet"`
	SnippetHighlightedWords []string `json:"snippet_highlighted_words"`
	Thumbnail               string   `json:"thumbnail,omitempty"` // omitempty 表示字段为空时省略
	Title                   string   `json:"title"`
}

type CallbackOutput struct {
	// Message is the message generated by the model.
	Message *schema.Message
	// Config is the config for the model.
	Config *Config
	// TokenUsage is the token usage of this request.
	TokenUsage *TokenUsage
	// Extra is the extra information for the callback.
	Extra map[string]any
}

type TokenUsage struct {
	// PromptTokens is the number of prompt tokens, including all the input tokens of this request.
	PromptTokens int
	// PromptTokenDetails is a breakdown of the prompt tokens.
	PromptTokenDetails any
	// CompletionTokens is the number of completion tokens.
	CompletionTokens int
	// TotalTokens is the total number of tokens.
	TotalTokens int
}

type Config struct {
	// Model is the model name.
	Model string
	// MaxTokens is the max number of tokens, if reached the max tokens, the model will stop generating, and mostly return an finish reason of "length".
	MaxTokens int
	// Temperature is the temperature, which controls the randomness of the model.
	Temperature float32
	// TopP is the top p, which controls the diversity of the model.
	TopP float32
	// Stop is the stop words, which controls the stopping condition of the model.
	Stop []string
}
