package tools

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

func TestPhotoSearchTool_InvokableRun(t1 *testing.T) {
	type fields struct {
		Query   string
		PerPage int
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
			name: "test",
			fields: fields{
				Query:   "computer",
				PerPage: 3,
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"query": "computer", "per_page": 3}`,
				in2:             []tool.Option{},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &PhotoSearchTool{
				Query:   tt.fields.Query,
				PerPage: tt.fields.PerPage,
			}
			got, err := t.InvokableRun(tt.args.ctx, tt.args.argumentsInJSON, tt.args.in2...)
			fmt.Println(got)
			fmt.Println(err)
		})
	}
}
