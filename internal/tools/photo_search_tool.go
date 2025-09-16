package tools

import (
	"agent/internal/consts"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type PhotoSearchTool struct {
	Query   string `json:"query"`
	PerPage int    `json:"per_page,omitempty"`
}

// PexelsPhoto 表示 Pexels API 返回的单张照片信息
type PexelsPhoto struct {
	ID              int    `json:"id"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	URL             string `json:"url"`
	Photographer    string `json:"photographer"`
	PhotographerURL string `json:"photographer_url"`
	AvgColor        string `json:"avg_color"`
	Src             struct {
		Original  string `json:"original"`
		Large2x   string `json:"large2x"`
		Large     string `json:"large"`
		Medium    string `json:"medium"`
		Small     string `json:"small"`
		Portrait  string `json:"portrait"`
		Landscape string `json:"landscape"`
		Tiny      string `json:"tiny"`
	} `json:"src"`
	Liked bool   `json:"liked"`
	Alt   string `json:"alt"`
}

// PexelsResponse 表示 Pexels API 的完整响应
type PexelsResponse struct {
	TotalResults int           `json:"total_results"`
	Page         int           `json:"page"`
	PerPage      int           `json:"per_page"`
	Photos       []PexelsPhoto `json:"photos"`
	NextPage     string        `json:"next_page"`
}

// PhotoResult 表示返回给用户的简化照片信息
type PhotoResult struct {
	ID           int    `json:"id"`
	Alt          string `json:"alt"`
	MediumURL    string `json:"medium_url"`
	Photographer string `json:"photographer"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

func NewPhotoSearchTool() *PhotoSearchTool {
	return &PhotoSearchTool{}
}

func (t *PhotoSearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "photo_search_tool",
		Desc: "Search for high-quality photos and return medium-sized image URLs",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     schema.String,
				Desc:     "Search query for photos (e.g., 'nature', 'city', 'people')",
				Required: true,
			},
			"per_page": {
				Type:     schema.Integer,
				Desc:     "Number of photos to return (1-80, default: 5)",
				Required: false,
			},
		}),
	}, nil
}

func (t *PhotoSearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	// 1. 反序列化参数
	var req PhotoSearchTool
	err := gjson.DecodeTo([]byte(argumentsInJSON), &req)
	if err != nil {
		return "", fmt.Errorf("failed to parse arguments: %v", err)
	}

	// 2. 参数验证和默认值设置
	if req.Query == "" {
		return "", fmt.Errorf("query parameter is required")
	}

	if req.PerPage <= 0 {
		req.PerPage = 5 // 默认返回5张图片
	}
	if req.PerPage > 80 {
		req.PerPage = 80 // Pexels API 限制最多80张
	}

	// 3. 获取 Pexels API Key
	apiKey := gconv.String(g.Cfg().MustGet(ctx, consts.PexelsApiKey))
	if apiKey == "" {
		return "", fmt.Errorf("Pexels API key not configured")
	}

	// 4. 构建请求URL
	baseURL := "https://api.pexels.com/v1/search"
	queryParams := url.Values{
		"query":    {req.Query},
		"per_page": {gconv.String(req.PerPage)},
	}
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// 5. 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 6. 设置Authorization头
	httpReq.Header.Set("Authorization", apiKey)
	httpReq.Header.Set("User-Agent", "PhotoSearchTool/1.0")

	// 7. 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 8. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 9. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// 10. 解析JSON响应
	var pexelsResp PexelsResponse
	err = gjson.DecodeTo(body, &pexelsResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// 11. 提取medium URL和相关信息
	var results []PhotoResult
	for _, photo := range pexelsResp.Photos {
		result := PhotoResult{
			ID:           photo.ID,
			Alt:          photo.Alt,
			MediumURL:    photo.Src.Medium,
			Photographer: photo.Photographer,
			Width:        photo.Width,
			Height:       photo.Height,
		}
		results = append(results, result)
	}

	// 12. 构建返回信息
	if len(results) == 0 {
		return fmt.Sprintf("No photos found for query: %s", req.Query), nil
	}

	// 13. 格式化返回结果
	var resultStrings []string
	resultStrings = append(resultStrings, fmt.Sprintf("Found %d photos for query '%s':", len(results), req.Query))

	for i, result := range results {
		resultStrings = append(resultStrings, fmt.Sprintf(
			"%d. %s\n   - Medium URL: %s\n   - Photographer: %s\n   - Size: %dx%d",
			i+1, result.Alt, result.MediumURL, result.Photographer, result.Width, result.Height,
		))
	}

	// 14. 同时返回JSON格式供程序使用
	jsonResult, err := gjson.EncodeString(results)
	if err != nil {
		return strings.Join(resultStrings, "\n"), nil
	}

	finalResult := strings.Join(resultStrings, "\n") + "\n\nJSON Data:\n" + jsonResult

	g.Log().Infof(ctx, "Photo search completed: query=%s, results=%d", req.Query, len(results))

	return finalResult, nil
}
