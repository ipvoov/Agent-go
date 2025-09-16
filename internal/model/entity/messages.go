// =================================================================================
// Code generated and maintained by GoFrame CLI tools. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Messages is the golang structure for table messages.
type Messages struct {
	Id          int64       `json:"id"          orm:"id"           description:""`     //
	SessionId   string      `json:"sessionId"   orm:"session_id"   description:"会话ID"` // 会话ID
	MessageType string      `json:"messageType" orm:"message_type" description:"消息类型"` // 消息类型
	Content     string      `json:"content"     orm:"content"      description:"消息内容"` // 消息内容
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"` // 创建时间
}
