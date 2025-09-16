// =================================================================================
// Code generated and maintained by GoFrame CLI tools. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Messages is the golang structure of table messages for DAO operations like Where/Data.
type Messages struct {
	g.Meta      `orm:"table:messages, do:true"`
	Id          any         //
	SessionId   any         // 会话ID
	MessageType any         // 消息类型
	Content     any         // 消息内容
	CreatedAt   *gtime.Time // 创建时间
}
