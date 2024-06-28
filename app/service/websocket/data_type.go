package websocket

import "goskeleton/app/model/msg"

// MsgPush 主动关爱消息格式
type MsgPush struct {
	Code int64 `json:"code"`
	Data []msg.MsgContentModel
}
