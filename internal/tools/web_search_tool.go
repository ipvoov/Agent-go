package tools

import (
	"agent/internal/consts"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type WebSearchTool struct {
	Q      string `json:"q"`
	Engine string `json:"engine"`
}

func NewWebSearchTool() *WebSearchTool {
	return &WebSearchTool{}
}

func (t *WebSearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "web_search_tool",
		Desc: "Search for information from Search Engine",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"q": {
				Type:     schema.String,
				Desc:     "search query keyword",
				Required: true},
			"engine": {
				Type:     schema.String,
				Desc:     "search engine name",
				Enum:     []string{"baidu", "google"},
				Required: true},
		}),
	}, nil
}

func (t *WebSearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	// 1. 反序列化 argumentsInJSON，处理 option 等
	var resp WebSearchTool
	err := gjson.DecodeTo([]byte(argumentsInJSON), &resp)
	if err != nil {
		return "", err
	}
	// 2. 处理业务逻辑
	baseURL := "https://www.searchapi.io/api/v1/search"
	queryParams := url.Values{
		"engine":  {resp.Engine},
		"q":       {resp.Q},
		"api_key": {gconv.String(g.Cfg().MustGet(ctx, consts.SearchApiKey))},
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
	req, _ := http.NewRequest("GET", fullURL, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// 3. 解析返回的JSON并提取前5条organic_results
	var searchResult map[string]interface{}
	err = gjson.DecodeTo(body, &searchResult)
	if err != nil {
		return "", err
	}

	// 提取organic_results
	organicResults := gconv.SliceAny(searchResult["organic_results"])

	var top5 []interface{}
	if len(organicResults) > 5 {
		top5 = organicResults[:5]
	} else {
		top5 = organicResults
	}

	result, err := gjson.EncodeString(top5)
	if err != nil {
		return "", err
	}

	return result, nil
}
