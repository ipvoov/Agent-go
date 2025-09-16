// =================================================================================
// Code generated and maintained by GoFrame CLI tools. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Sessions is the golang structure for table sessions.
type Sessions struct {
	Id        int64       `json:"id"        orm:"id"         description:""`       //
	SessionId string      `json:"sessionId" orm:"session_id" description:"会话唯一标识"` // 会话唯一标识
	Title     string      `json:"title"     orm:"title"      description:"会话标题"`   // 会话标题
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`   // 更新时间
}
