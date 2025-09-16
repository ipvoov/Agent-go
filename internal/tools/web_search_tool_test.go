package tools

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

//func TestWebSearchTool(t *testing.T) {
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{name: "01", args: args{ctx: context.Background()}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			WebSearchTool(tt.args.ctx)
//		})
//	}
//}

func TestWebSearchTool_InvokableRun(t1 *testing.T) {
	type fields struct {
		Q      string
		Engine string
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
			name: "01",
			fields: fields{
				Q:      "丁真",
				Engine: "google",
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"q":"丁真","engine":"google"}`,
				in2:             []tool.Option{},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &WebSearchTool{
				Q:      tt.fields.Q,
				Engine: tt.fields.Engine,
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
