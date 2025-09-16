package tools

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

func TestFileOperationTool_InvokableRun(t1 *testing.T) {
	type fields struct {
		Operation string
		FilePath  string
		Content   string
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
			name: "TestFileOperationTool_InvokableRun",
			fields: fields{
				Operation: "write",
				FilePath:  "resource/txt/test.txt",
				Content:   "test is a test",
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"operation": "create", "file_path": "resource/txt/test.txt", "content": "test is a test"}`,
				in2:             []tool.Option{},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &FileOperationTool{
				Operation: tt.fields.Operation,
				FilePath:  tt.fields.FilePath,
				Content:   tt.fields.Content,
			}
			got, err := t.InvokableRun(tt.args.ctx, tt.args.argumentsInJSON, tt.args.in2...)
			fmt.Println(got)
			fmt.Println(err)
		})
	}
}
