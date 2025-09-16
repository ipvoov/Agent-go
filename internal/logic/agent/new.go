package agent

import (
	v1 "agent/api/agent/v1"
	"agent/internal/consts"
	"agent/internal/service"
	"context"
	"fmt"
	"log"

	askembedding "github.com/cloudwego/eino-ext/components/embedding/ark"
	askmodel "github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	clie "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func GetGdMapMCPTool(ctx context.Context) tool.BaseTool {
	url := "https://mcp.amap.com/sse?key=" + g.Cfg().MustGet(ctx, consts.McpApiKey).String()
	cli, err := clie.NewSSEMCPClient(url)
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatal(err)
	}

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli, ToolNameList: []string{"maps_around_search"}})
	if err != nil {
		log.Fatal(err)
	}

	return tools[0]
}

func NewChatModel(ctx context.Context) *askmodel.ChatModel {
	chatModel, err := askmodel.NewChatModel(ctx, &askmodel.ChatModelConfig{
		APIKey: g.Cfg().MustGet(ctx, consts.ApiKey).String(),
		Model:  g.Cfg().MustGet(ctx, consts.Model).String(),
	})
	if err != nil {
		panic(err)
	}
	return chatModel
}

func NewMilVusRetriever(ctx context.Context) *milvus.Retriever {
	cli, err := client.NewClient(ctx, client.Config{
		Address: g.Cfg().MustGet(ctx, consts.MilvusAddr).String(),
	})
	if err != nil {
		panic(err)
	}

	emb, err := askembedding.NewEmbedder(ctx, &askembedding.EmbeddingConfig{
		APIKey: g.Cfg().MustGet(ctx, consts.ApiKey).String(),
		Model:  g.Cfg().MustGet(ctx, consts.EmbModel).String(),
	})
	if err != nil {
		panic(err)
	}

	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  "test",
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:           2,
		ScoreThreshold: 0.97, // 提高到 97%，非常严格
		Embedding:      emb,
	})
	if err != nil {
		panic(err)
	}
	return retriever
}

func Template(ctx context.Context, in *v1.ChatStreamReq) (*prompt.DefaultChatTemplate, string) {
	// 格式化模板
	// 创建模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`you are an {role}.
You have in-depth research on the psychology of love and are good at analyzing emotional dynamics,
providing practical love advice, and helping people build and maintain healthy relationships.
Provide users with professional love guidance to help them solve love problems, improve love skills, and promote the healthy development of emotional relationships.
{example}
Answer based on the following content.
If you need more information, you can use the search tool, but you can only use it once.
`),
		schema.MessagesPlaceholder("history_key", false),
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)

	example := "- Answer based on the following content:"
	rag := NewMilVusRetriever(ctx)
	result, err := rag.Retrieve(ctx, in.Query)
	// 添加额外的相关性检查
	if err == nil && len(result) > 0 && service.Agent().IsRelevant(in.Query, result) {
		for i, doc := range result {
			// 限制每个文档的长度，避免 token 过多
			content := doc.Content
			if len(content) > 1000 {
				content = content[:1000] + "..."
			}
			example += fmt.Sprintf("\n example %d:\n%s\n", i+1, content)
		}
	}

	return template, example
}

func AgentTemplate(ctx context.Context, in *v1.ChatStreamReq) []*schema.Message {
	// 格式化模板
	// 创建模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`you are an {role}.
You have natural language processing knowledge, interactive design capabilities,
logical analysis skills, and a deep understanding of artificial intelligence agents,
and List your thinking steps.
`),
		schema.MessagesPlaceholder("history_key", false),
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)

	sessionMessages := service.Agent().GetSessionBySessionID(ctx, in.SessionID)
	variables := map[string]any{
		"role":        "Artificial intelligence interaction experts and intelligent agent consultants",
		"task":        in.Query,
		"history_key": sessionMessages,
	}
	messages, _ := template.Format(ctx, variables)

	return messages
}
