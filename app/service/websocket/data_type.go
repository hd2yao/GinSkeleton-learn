package websocket

import "goskeleton/app/model/msg"

// MsgPush 主动关爱消息格式
type MsgPush struct {
	Code int64 `json:"code"`
	Data []msg.MsgContentModel
}

// ReturnClientMsg 接受到 ws 客户端发送过来的消息
type ReturnClientMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Type               int64  `json:"type"`
		UserId             int64  `json:"user_id"`
		RoomId             int64  `json:"room_id"`
		HandFree           int    `json:"hands_free"`
		DeviceType         int    `json:"device_type"`
		UserIdCallerTitle  string `json:"user_id_caller_title"`
		UserIdCallerAvatar string `json:"user_id_caller_avatar"`
	} `json:"data"`
}
