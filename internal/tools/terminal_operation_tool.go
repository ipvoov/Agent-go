package tools

import (
	"agent/internal/consts"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

type TerminalOperationTool struct {
	Command   string `json:"command"`
	Directory string `json:"directory,omitempty"`
}

func NewTerminalOperationTool() *TerminalOperationTool {
	return &TerminalOperationTool{}
}

func (t *TerminalOperationTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "terminal_operation_tool",
		Desc: "Execute terminal/shell commands and return the output",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"command": {
				Type:     schema.String,
				Desc:     "The command to execute in terminal",
				Required: true,
			},
			"directory": {
				Type:     schema.String,
				Desc:     "Working directory for the command (optional, defaults to current directory)",
				Required: false,
			},
		}),
	}, nil
}

func (t *TerminalOperationTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	// 1. 反序列化 argumentsInJSON
	var req TerminalOperationTool
	err := gjson.DecodeTo([]byte(argumentsInJSON), &req)
	if err != nil {
		return "", fmt.Errorf("failed to parse arguments: %v", err)
	}

	// 2. 参数验证和默认值设置
	if req.Command == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	if req.Directory == "" {
		req.Directory, _ = os.Getwd() // 默认当前目录
	}

	lowerCmd := strings.ToLower(req.Command)
	if gstr.InArray(consts.DangerousCommands, lowerCmd) {
		return "", fmt.Errorf("dangerous command detected and blocked: %s", lowerCmd)
	}

	// 4. 创建命令执行上下文
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(30)*time.Second)
	defer cancel()

	// 5. 根据操作系统选择shell
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(cmdCtx, "cmd", "/C", req.Command)
	default:
		cmd = exec.CommandContext(cmdCtx, "sh", "-c", req.Command)
	}

	// 6. 设置工作目录
	cmd.Dir = req.Directory

	// 7. 执行命令并获取输出
	g.Log().Infof(ctx, "Executing command: %s in directory: %s", req.Command, req.Directory)

	output, err := cmd.CombinedOutput()

	g.Log().Infof(ctx, "Command output: %s", string(output))

	return string(output), err

	// 8. 构建返回结果
	//result := map[string]any{
	//	"command":   req.Command,
	//	"directory": req.Directory,
	//	"success":   err == nil,
	//	"output":    outputStr,
	//	"timestamp": gtime.Now().Format("Y-m-d H:i:s"),
	//}
	//
	//if err != nil {
	//	result["error"] = err.Error()
	//	g.Log().Warningf(ctx, "Command execution failed: %v", err)
	//} else {
	//	g.Log().Infof(ctx, "Command executed successfully")
	//}
	//
	//// 9. 返回JSON格式的结果
	//resultJSON, err := gjson.EncodeString(result)
	//if err != nil {
	//	return "", fmt.Errorf("failed to encode result: %v", err)
	//}
	//
	//return resultJSON, nil
}
