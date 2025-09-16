package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	_ "agent/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"agent/internal/cmd"
)

func main() {
	// -------------初始化 coze loop 客户端----------
	// 设置相关环境变量
	// COZELOOP_WORKSPACE_ID=your workspace id
	// COZELOOP_API_TOKEN=your token
	//client, err := cozeloop.NewClient()
	//if err != nil {
	//	panic(err)
	//}
	//defer client.Close(ctx)
	//// 在服务 init 时 once 调用
	//handler := ccb.NewLoopHandler(client)
	//callbacks.AppendGlobalHandlers(handler)

	cmd.Main.Run(gctx.GetInitCtx())
}
