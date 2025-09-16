package consts

const (
	ApiKey       = "ai.apiKey"
	Model        = "ai.model"
	EmbModel     = "ai.embModel"
	MilvusAddr   = "ai.milvusAddr"
	SearchApiKey = "ai.SearchApiKey"
	PexelsApiKey = "ai.pexelsApiKey"
	McpApiKey    = "ai.mcpApiKey"

	System    = "system"
	User      = "user"
	Assistant = "assistant"

	Rag           = "rag"
	Ai            = "ai"
	Tools         = "tools"
	DefaultSysMsg = "你是豆包，是字节跳动研发的人工智能助手，你可以回答用户的问题"
)

var (
	//DangerousCommands 安全检查 - 禁止危险命令
	DangerousCommands = []string{
		"rm -rf", "sudo rm", "format", "del /f", "rmdir /s",
		"shutdown", "reboot", "halt", "poweroff",
		"dd if=", "mkfs", "fdisk",
	}
)
