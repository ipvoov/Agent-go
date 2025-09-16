package tools

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

func TestTerminalOperationTool_InvokableRun(t1 *testing.T) {
	type fields struct {
		Command   string
		Directory string
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
			name: "TestTerminalOperationTool_InvokableRun",
			fields: fields{
				Command:   "ls",
				Directory: "/Users/offves/GolandProjects/agent",
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"command": "ls", "directory": "/Users/offves/GolandProjects/agent"}`,
				in2:             []tool.Option{},
			},
			want:    `{"command":"ls","directory":"/Users/offves/GolandProjects/agent","success":true,"output":"consts\ninternal\ntest\n","timestamp":"2023-08-10 15:04:05"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TerminalOperationTool{
				Command:   tt.fields.Command,
				Directory: tt.fields.Directory,
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
