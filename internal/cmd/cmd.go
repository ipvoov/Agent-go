package cmd

import (
	"context"

	"agent/internal/controller/agent"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// -------------初始化 http 服务----------
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(
					agent.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
