package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
)

type ResourceDownloadTool struct {
	URL      string `json:"url"`
	Filename string `json:"filename,omitempty"`
}

func NewResourceDownloadTool() *ResourceDownloadTool {
	return &ResourceDownloadTool{}
}

func (t *ResourceDownloadTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "resource_download_tool",
		Desc: "Download files from URL and save to resource/download directory",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"url": {
				Type:     schema.String,
				Desc:     "The URL of the resource to download",
				Required: true,
			},
			"filename": {
				Type:     schema.String,
				Desc:     "Custom filename for the downloaded file (optional, will auto-detect from URL if not provided)",
				Required: false,
			},
		}),
	}, nil
}

func (t *ResourceDownloadTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	// 1. 反序列化参数
	var req ResourceDownloadTool
	err := gjson.DecodeTo([]byte(argumentsInJSON), &req)
	if err != nil {
		return "", fmt.Errorf("failed to parse arguments: %v", err)
	}

	// 2. 参数验证和默认值设置
	if req.URL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// 验证URL格式
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("only HTTP and HTTPS URLs are supported")
	}

	// 3. 确定文件名
	filename := req.Filename
	if filename == "" {
		// 从URL中提取文件名
		filename = filepath.Base(parsedURL.Path)
		if filename == "." || filename == "/" {
			// 如果无法从URL提取文件名，使用时间戳
			filename = fmt.Sprintf("download_%d", time.Now().Unix())
		}
	}

	// 安全检查：确保文件名不包含路径遍历字符
	filename = filepath.Base(filename)
	if filename == "" || filename == "." || filename == ".." {
		filename = fmt.Sprintf("download_%d", time.Now().Unix())
	}

	// 4. 创建下载目录
	downloadDir := "resource/download"
	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create download directory: %v", err)
	}

	// 构建完整文件路径
	filePath := filepath.Join(downloadDir, filename)

	// 检查文件是否已存在，如果存在则添加序号
	originalPath := filePath
	counter := 1
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(originalPath)
		nameWithoutExt := strings.TrimSuffix(filepath.Base(originalPath), ext)
		filePath = filepath.Join(downloadDir, fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext))
		counter++
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download from URL: %v", err)
	}
	defer resp.Body.Close()

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	//下载文件内容
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		// 如果下载失败，删除部分下载的文件
		os.Remove(filePath)
		return "", fmt.Errorf("failed to download file content: %v", err)
	}

	return fmt.Sprintf("Resource downloaded successfully to %s", filePath), nil
}
