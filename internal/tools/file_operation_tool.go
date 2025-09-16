package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

type FileOperationTool struct {
	Operation string `json:"operation"`
	FilePath  string `json:"file_path"`
	Content   string `json:"content,omitempty"`
}

func NewFileOperationTool() *FileOperationTool {
	return &FileOperationTool{}
}

func (t *FileOperationTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "file_operation_tool",
		Desc: "Perform basic file operations (read, write, create, delete) on files in the current directory and its subdirectories",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"operation": {
				Type:     schema.String,
				Desc:     "File operation to perform",
				Enum:     []string{"read", "write", "create", "delete"},
				Required: true,
			},
			"file_path": {
				Type:     schema.String,
				Desc:     "File name or relative path within current directory (e.g., 'file.txt', 'resource/pdf/doc.pdf')",
				Required: true,
			},
			"content": {
				Type:     schema.String,
				Desc:     "Content to write (for write/create operations)",
				Required: false,
			},
		}),
	}, nil
}

func (t *FileOperationTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	// 1. 反序列化参数
	var req FileOperationTool
	err := gjson.DecodeTo([]byte(argumentsInJSON), &req)
	if err != nil {
		return "", fmt.Errorf("failed to parse arguments: %v", err)
	}

	// 2. 参数验证
	if req.Operation == "" {
		return "", fmt.Errorf("operation parameter is required")
	}
	if req.FilePath == "" {
		return "", fmt.Errorf("file_path parameter is required")
	}

	// 3. 安全检查 - 确保只能操作当前目录的文件
	if err := t.validateCurrentDirPath(req.FilePath); err != nil {
		return "", err
	}

	// 4. 根据操作类型执行相应的文件操作
	switch strings.ToLower(req.Operation) {
	case "read":
		return t.readFile(ctx, req.FilePath)
	case "write":
		return t.writeFile(ctx, req.FilePath, req.Content)
	case "create":
		return t.createFile(ctx, req.FilePath, req.Content)
	case "delete":
		return t.deleteFile(ctx, req.FilePath)
	default:
		return "", fmt.Errorf("unsupported operation: %s", req.Operation)
	}
}

// validateCurrentDirPath 验证路径只能在当前目录及其子目录内
func (t *FileOperationTool) validateCurrentDirPath(path string) error {
	// 清理路径，防止路径遍历攻击
	cleanPath := filepath.Clean(path)

	// 检查是否包含路径遍历字符
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	// 检查是否为绝对路径
	if filepath.IsAbs(cleanPath) {
		return fmt.Errorf("absolute paths not allowed, only relative paths within current directory: %s", path)
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// 获取目标文件的绝对路径
	absPath := filepath.Join(currentDir, cleanPath)

	// 确保目标路径在当前目录内（包括子目录）
	// 添加路径分隔符确保是真正的子目录或当前目录
	currentDirWithSep := currentDir + string(filepath.Separator)
	absPathWithSep := absPath + string(filepath.Separator)

	// 检查是否在当前目录或其子目录中
	if !strings.HasPrefix(absPathWithSep, currentDirWithSep) && absPath != currentDir {
		return fmt.Errorf("file must be within current directory or its subdirectories: %s", path)
	}

	return nil
}

// readFile 读取文件内容
func (t *FileOperationTool) readFile(ctx context.Context, filePath string) (string, error) {
	if !gfile.Exists(filePath) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	if gfile.IsDir(filePath) {
		return "", fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	content := gfile.GetContents(filePath)
	if content == "" {
		// 检查是否是因为文件为空还是读取失败
		if stat, err := os.Stat(filePath); err == nil && stat.Size() == 0 {
			return "File is empty", nil
		}
		return "", fmt.Errorf("failed to read file or file is empty: %s", filePath)
	}

	g.Log().Infof(ctx, "Successfully read file: %s (%d bytes)", filePath, len(content))
	return fmt.Sprintf("File content of %s:\n%s", filePath, content), nil
}

// writeFile 写入文件内容
func (t *FileOperationTool) writeFile(ctx context.Context, filePath, content string) (string, error) {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if dir != "." && !gfile.Exists(dir) {
		err := gfile.Mkdir(dir)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		g.Log().Infof(ctx, "Created directory: %s", dir)
	}

	err := gfile.PutContents(filePath, content)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	g.Log().Infof(ctx, "Successfully wrote to file: %s (%d bytes)", filePath, len(content))
	return fmt.Sprintf("Successfully wrote %d bytes to file: %s", len(content), filePath), nil
}

// createFile 创建新文件
func (t *FileOperationTool) createFile(ctx context.Context, filePath, content string) (string, error) {
	if gfile.Exists(filePath) {
		return "", fmt.Errorf("file already exists: %s", filePath)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if dir != "." && !gfile.Exists(dir) {
		err := gfile.Mkdir(dir)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		g.Log().Infof(ctx, "Created directory: %s", dir)
	}

	err := gfile.PutContents(filePath, content)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	g.Log().Infof(ctx, "Successfully created file: %s (%d bytes)", filePath, len(content))
	return fmt.Sprintf("Successfully created file: %s with %d bytes", filePath, len(content)), nil
}

// deleteFile 删除文件
func (t *FileOperationTool) deleteFile(ctx context.Context, filePath string) (string, error) {
	if !gfile.Exists(filePath) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	if gfile.IsDir(filePath) {
		return "", fmt.Errorf("cannot delete directories, only files are supported: %s", filePath)
	}

	err := os.Remove(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to delete file %s: %v", filePath, err)
	}

	g.Log().Infof(ctx, "Successfully deleted file: %s", filePath)
	return fmt.Sprintf("Successfully deleted file: %s", filePath), nil
}
