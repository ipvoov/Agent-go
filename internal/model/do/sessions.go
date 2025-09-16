// =================================================================================
// Code generated and maintained by GoFrame CLI tools. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Sessions is the golang structure of table sessions for DAO operations like Where/Data.
type Sessions struct {
	g.Meta    `orm:"table:sessions, do:true"`
	Id        any         //
	SessionId any         // 会话唯一标识
	Title     any         // 会话标题
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
