package tools

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

func TestPDFGenerationTool_InvokableRun(t1 *testing.T) {
	type fields struct {
		Filename string
		Content  string
		Title    string
		Author   string
		Subject  string
	}
	type args struct {
		ctx             context.Context
		argumentsInJSON string
		in2             []tool.Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestPDFGenerationTool_InvokableRun_HTML",
			fields: fields{
				Filename: "test.pdf",
				Content:  "<h1>主标题</h1><h2>副标题</h2><p>这是正文内容，支持<strong>粗体</strong>和<em>斜体</em>。</p><table><tr><th>列1</th><th>列2</th></tr><tr><td>数据1</td><td>数据2</td></tr></table>",
				Title:    "HTML文档",
				Author:   "Ai助手",
				Subject:  "HTML文档生成测试",
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"filename": "test.pdf", "content": "<h1>主标题</h1><h2>副标题</h2><p>这是正文内容，支持<strong>粗体</strong>和<em>斜体</em>。</p><table><tr><th>列1</th><th>列2</th></tr><tr><td>数据1</td><td>数据2</td></tr></table>", "title": "HTML文档", "author": "Ai助手", "subject": "HTML文档生成测试"}`,
				in2:             []tool.Option{},
			},
			want:    "PDF generation to resource/pdf/test.pdf",
			wantErr: false,
		},
		{
			name: "TestPDFGenerationTool_InvokableRun",
			fields: fields{
				Filename: "test01.pdf",
				Content:  "这是第一段文字。\n\n这是第二段文字。\n包含多行内容。",
				Title:    "HTML文档",
				Author:   "Ai助手",
				Subject:  "HTML文档生成测试",
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"filename": "test01.pdf", "content": "这是第一段文字。\n\n这是第二段文字。\n包含多行内容。", "title": "HTML文档", "author": "Ai助手", "subject": "HTML文档生成测试"}`,
				in2:             []tool.Option{},
			},
			want:    "PDF generation to resource/pdf/test01.pdf",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &PDFGenerationTool{
				Filename: tt.fields.Filename,
				Content:  tt.fields.Content,
				Title:    tt.fields.Title,
				Author:   tt.fields.Author,
				Subject:  tt.fields.Subject,
			}
			got, err := t.InvokableRun(tt.args.ctx, tt.args.argumentsInJSON, tt.args.in2...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("InvokableRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("InvokableRun() got = %v, want %v", got, tt.want)
			}
		})
	}
}
