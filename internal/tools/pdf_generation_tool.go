package tools

import (
	"context"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type PDFGenerationTool struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Subject  string `json:"subject,omitempty"`
}

func NewPDFGenerationTool() *PDFGenerationTool {
	return &PDFGenerationTool{}
}

func (t *PDFGenerationTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "pdf_generation_tool",
		Desc: "Generate PDF files from HTML content and save to resource/pdf directory",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"filename": {
				Type:     schema.String,
				Desc:     "The filename for the generated PDF (without .pdf extension)",
				Required: true,
			},
			"content": {
				Type:     schema.String,
				Desc:     "The HTML content to be included in the PDF. Plain text will be automatically wrapped in HTML tags.",
				Required: true,
			},
			"title": {
				Type:     schema.String,
				Desc:     "PDF document title (optional)",
				Required: false,
			},
			"author": {
				Type:     schema.String,
				Desc:     "PDF document author (optional)",
				Required: false,
			},
			"subject": {
				Type:     schema.String,
				Desc:     "PDF document subject (optional)",
				Required: false,
			},
		}),
	}, nil
}

func (t *PDFGenerationTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	// 1. 反序列化参数
	var req PDFGenerationTool
	err := gjson.DecodeTo([]byte(argumentsInJSON), &req)
	if err != nil {
		return "Error generation PDF: failed to parse arguments: " + err.Error(), nil
	}

	// 2. 参数验证
	if req.Filename == "" {
		return "Error generation PDF: filename cannot be empty", nil
	}

	if req.Content == "" {
		return "Error generation PDF: content cannot be empty", nil
	}

	// 3. 处理文件名
	filename := strings.TrimSpace(req.Filename)
	filename = sanitizeFilename(filename)
	if filename == "" {
		filename = fmt.Sprintf("document_%d", time.Now().Unix())
	}

	// 确保文件名以.pdf结尾
	if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		filename += ".pdf"
	}

	// 4. 创建PDF目录
	pdfDir := "resource/pdf"
	err = os.MkdirAll(pdfDir, 0755)
	if err != nil {
		return "Error generation PDF: failed to create PDF directory: " + err.Error(), nil
	}

	// 5. 构建完整文件路径
	filePath := filepath.Join(pdfDir, filename)

	// 检查文件是否已存在，如果存在则添加序号
	originalPath := filePath
	counter := 1
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(originalPath)
		nameWithoutExt := strings.TrimSuffix(filepath.Base(originalPath), ext)
		filePath = filepath.Join(pdfDir, fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext))
		counter++
	}

	// 6. 生成PDF
	return generatePDFFromHTML(ctx, req, filePath)
}

// generatePDFFromHTML 使用chromedp从HTML生成PDF
func generatePDFFromHTML(ctx context.Context, req PDFGenerationTool, filePath string) (string, error) {
	// 1. 构建完整的HTML内容
	htmlContent := buildHTMLContent(req)

	// 调试：输出HTML内容长度
	g.Log().Infof(ctx, "Generated HTML content length: %d", len(htmlContent))

	// 2. 创建chromedp上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.WindowSize(1200, 800),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 3. 设置超时
	chromeCtx, cancel = context.WithTimeout(chromeCtx, 30*time.Second)
	defer cancel()

	// 4. 生成PDF
	g.Log().Infof(ctx, "Generating PDF from HTML using Chrome: %s", filePath)

	// 创建临时HTML文件而不是使用data URL
	tempDir := os.TempDir()
	tempHTMLFile := filepath.Join(tempDir, fmt.Sprintf("temp_pdf_%d.html", time.Now().UnixNano()))
	defer os.Remove(tempHTMLFile) // 清理临时文件

	// 写入HTML内容到临时文件
	err := os.WriteFile(tempHTMLFile, []byte(htmlContent), 0644)
	if err != nil {
		return "Error generation PDF: failed to create temporary HTML file: " + err.Error(), nil
	}

	// 使用file:// URL
	fileURL := "file://" + tempHTMLFile

	var pdfBuffer []byte
	err = chromedp.Run(chromeCtx,
		// 导航到临时HTML文件
		chromedp.Navigate(fileURL),
		// 等待页面加载完成
		chromedp.WaitReady("body"),
		// 等待内容渲染
		chromedp.Sleep(3*time.Second),
		// 生成PDF
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().
				WithPaperWidth(8.27).  // A4宽度（英寸）
				WithPaperHeight(11.7). // A4高度（英寸）
				WithMarginTop(0.5).    // 边距
				WithMarginBottom(0.5).
				WithMarginLeft(0.5).
				WithMarginRight(0.5).
				WithPrintBackground(true).
				WithDisplayHeaderFooter(false).
				WithScale(1.0).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		g.Log().Errorf(ctx, "Chrome PDF generation error: %v", err)
		return "Error generation PDF: " + err.Error(), nil
	}

	// 检查PDF缓冲区大小
	g.Log().Infof(ctx, "Generated PDF buffer size: %d bytes", len(pdfBuffer))

	if len(pdfBuffer) == 0 {
		return "Error generation PDF: generated PDF is empty", nil
	}

	// 5. 写入PDF文件
	err = os.WriteFile(filePath, pdfBuffer, 0644)
	if err != nil {
		return "Error generation PDF: failed to write PDF file: " + err.Error(), nil
	}

	g.Log().Infof(ctx, "PDF generated successfully from HTML: %s", filepath.Base(filePath))

	// 6. 返回成功信息
	return "PDF generation successfully to " + filePath, nil
}

// buildHTMLContent 构建完整的HTML内容
func buildHTMLContent(req PDFGenerationTool) string {
	// 检查内容是否已经是HTML格式
	content := req.Content
	isHTML := strings.Contains(strings.ToLower(content), "<html") ||
		strings.Contains(strings.ToLower(content), "<!doctype") ||
		strings.Contains(content, "<h1>") ||
		strings.Contains(content, "<h2>") ||
		strings.Contains(content, "<p>") ||
		strings.Contains(content, "<div>")

	// 如果不是HTML格式，将纯文本转换为HTML段落
	if !isHTML {
		lines := strings.Split(content, "\n")
		var htmlLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				htmlLines = append(htmlLines, "<br>")
			} else {
				htmlLines = append(htmlLines, "<p>"+html.EscapeString(line)+"</p>")
			}
		}
		content = strings.Join(htmlLines, "\n")
	}

	// 构建完整的HTML文档
	htmlDoc := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + html.EscapeString(req.Title) + `</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: "PingFang SC", "Microsoft YaHei", "Hiragino Sans GB", "SimSun", Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: white;
            padding: 20px;
            font-size: 14px;
        }
        
        .document {
            max-width: 800px;
            margin: 0 auto;
            background: white;
        }
        
        .header {
            margin-bottom: 40px;
            border-bottom: 3px solid #3498db;
            padding-bottom: 25px;
            text-align: center;
        }
        
        .title {
            font-size: 28px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 15px;
            letter-spacing: 1px;
        }
        
        .meta-info {
            font-size: 13px;
            color: #7f8c8d;
            background-color: #f8f9fa;
            padding: 10px 20px;
            border-radius: 5px;
            display: inline-block;
        }
        
        .content {
            line-height: 1.8;
            text-align: justify;
            font-size: 15px;
        }
        
        .content h1 {
            font-size: 22px;
            color: #2c3e50;
            margin: 30px 0 20px 0;
            border-bottom: 2px solid #3498db;
            padding-bottom: 8px;
        }
        
        .content h2 {
            font-size: 19px;
            color: #34495e;
            margin: 25px 0 15px 0;
            border-left: 4px solid #3498db;
            padding-left: 10px;
        }
        
        .content h3 {
            font-size: 17px;
            color: #34495e;
            margin: 20px 0 12px 0;
        }
        
        .content p {
            margin-bottom: 16px;
            text-indent: 2em;
        }
        
        .content ul, .content ol {
            margin: 10px 0 10px 20px;
        }
        
        .content li {
            margin-bottom: 5px;
        }
        
        .content strong {
            font-weight: bold;
            color: #2c3e50;
        }
        
        .content em {
            font-style: italic;
        }
        
        .content blockquote {
            border-left: 4px solid #3498db;
            padding-left: 15px;
            margin: 15px 0;
            color: #666;
            font-style: italic;
        }
        
        .content table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
        }
        
        .content th, .content td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        
        .content th {
            background-color: #f5f5f5;
            font-weight: bold;
        }
        
        .footer {
            margin-top: 50px;
            padding-top: 20px;
            border-top: 1px solid #e0e0e0;
            font-size: 11px;
            color: #bdc3c7;
            text-align: center;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="document">
        <div class="header">
            <div class="title">` + html.EscapeString(req.Title) + `</div>
            <div class="meta-info">`

	if req.Author != "" {
		htmlDoc += `作者: ` + html.EscapeString(req.Author)
		if req.Subject != "" {
			htmlDoc += ` | `
		}
	}
	if req.Subject != "" {
		htmlDoc += `主题: ` + html.EscapeString(req.Subject)
	}

	htmlDoc += `<br>生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `
            </div>
        </div>
        
        <div class="content">
            ` + content + `
        </div>
        
        <div class="footer">
            Generated by PDF Generation Tool
        </div>
    </div>
</body>
</html>`

	return htmlDoc
}

// sanitizeFilename 清理文件名中的不安全字符
func sanitizeFilename(filename string) string {
	// 移除或替换不安全的字符
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range unsafe {
		filename = strings.ReplaceAll(filename, char, "_")
	}

	// 移除前后空格和点
	filename = strings.Trim(filename, " .")

	return filename
}
