package websocket

import "goskeleton/app/model/msg"

// MsgPush 主动关爱消息格式
type MsgPush struct {
	Code int64 `json:"code"`
	Data []msg.MsgContentModel
}

// ReturnClientMsg ws 返回给客户端的消息
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

// ReturnFiled 返回给客户端失败的信息
type ReturnFiled struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

// ReceivedWsClientMsg 接受到 ws 客户端发送过来的消息
type ReceivedWsClientMsg struct {
	Type   int64 `json:"type"`
	UserId int64 `json:"user_id"`
	IsCall int64 `json:"is_call"`
	Code   int64 `json:"code"`
}
