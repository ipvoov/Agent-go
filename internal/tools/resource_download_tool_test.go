package tools

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

func TestResourceDownloadTool_InvokableRun(t1 *testing.T) {
	type fields struct {
		URL      string
		Filename string
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
			name: "TestResourceDownloadTool_InvokableRun_Success",
			fields: fields{
				URL:      "https://picsum.photos/800/600",
				Filename: "image.jpg",
			},
			args: args{
				ctx:             context.Background(),
				argumentsInJSON: `{"url": "https://picsum.photos/800/600", "filename": "image.jpg"}`,
				in2:             []tool.Option{},
			},
			want:    "Resource downloaded successfully to image.jpg",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &ResourceDownloadTool{
				URL:      tt.fields.URL,
				Filename: tt.fields.Filename,
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
