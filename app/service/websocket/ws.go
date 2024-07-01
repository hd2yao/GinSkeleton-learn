package websocket

import (
	"encoding/json"
	"fmt"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/call_log"
	"goskeleton/app/model/friend_user"
	"goskeleton/app/model/home_user"
	"goskeleton/app/model/oldster_user"
	"goskeleton/app/model/room"
	"goskeleton/app/utils/websocket/core"
)

/**
websocket模块相关事件执行顺序：
1.onOpen
2.OnMessage
3.OnError
4.OnClose
*/

type Ws struct {
	WsClient *core.Client
}

// OnOpen 事件函数
func (w *Ws) OnOpen(context *gin.Context) (*Ws, bool) {
	if client, ok := (&core.Client{}).OnOpen(context); ok {

		//token := context.GetString(consts.ValidatorPrefix + "token")
		//variable.ZapLog.Info("获取到的客户端上线时携带的唯一标记值：", zap.String("token", token))

		// 成功上线以后，开发者可以基于客户端上线时携带的唯一参数(这里用token键表示)
		// 在数据库查询更多的其他字段信息，直接追加在 Client 结构体上，方便后续使用
		//client.ClientMoreParams.UserParams1 = "123"
		//client.ClientMoreParams.UserParams2 = "456"
		//fmt.Printf("最终每一个客户端(client) 已有的参数：%+v\n", client)

		client.UserId = int64(context.GetFloat64(consts.ValidatorPrefix + "user_id"))
		client.HomeId = int64(context.GetFloat64(consts.ValidatorPrefix + "home_id"))
		client.UserType = context.GetString(consts.ValidatorPrefix + "user_type")
		// 区域用于主动关爱推送消息
		userInfo := oldster_user.CreateOldsterUserModelFactory("").GetById(client.UserId)
		client.CityId = int64(userInfo.FkProvinceCityId)

		// 触发 onOpen 后，推送全部的信息
		HandleMsg(client, true)

		w.WsClient = client
		variable.ZapLog.Info("用户上线:ID:" + strconv.Itoa(int(w.WsClient.HomeId)) + "类型：" + w.WsClient.UserType)
		go w.WsClient.Heartbeat() // 一旦握手+协议升级成功，就为每一个连接开启一个自动化的隐式心跳检测包

		// 上线的如果是小程序，则判断是否在呼叫过程中，如果在呼叫过程中，则发送消息
		// 此时上线的小程序用户是被呼叫方
		if client.UserType == "mobile" {
			// 查询通话状态
			callId := home_user.CreateHomeModelFactory("").GetCallId(int(client.UserId))
			if callId != 0 {
				// 将呼叫方的信息返回给当前上线的小程序用户
				returnCalled := ReturnClientMsg{}
				returnCalled.Code = 200
				returnCalled.Msg = "success"
				returnCalled.Data.Type = 1
				returnCalled.Data.UserId = int64(callId)
				// 呼叫方创建的视频通话房间
				returnCalled.Data.RoomId = room.CreateRoomModelFactory("").GetRoomId(callId)

				// 查询是否给呼叫方设置备注以及呼叫方的相关信息
				// GetByUserIdAndFriendId() 方法中会先初始化，因此当没有查询到结果，返回的也都是数据类型的默认值
				callFriendData := friend_user.CreateFriendUserModelFactory("").GetByUserIdAndFriendId(int(w.WsClient.HomeId), callId)
				userInfo := home_user.CreateHomeModelFactory("").GetHomeUser(int64(callId))
				returnCalled.Data.DeviceType = userInfo.DeviceType
				returnCalled.Data.UserIdCallerAvatar = userInfo.Avatar
				// 也就是说 当没有好友信息时，callFriendData 是一个空的 FriendUserModel{} 而不是 nil
				returnCalled.Data.HandFree = callFriendData.HandsFree
				if callFriendData.NickName != "" {
					returnCalled.Data.UserIdCallerTitle = callFriendData.NickName
				} else {
					returnCalled.Data.UserIdCallerTitle = userInfo.Title
				}
				returnCalledStr, _ := json.Marshal(returnCalled)
				if err := w.WsClient.SendMessage(1, string(returnCalledStr)); err != nil {
					variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
				} else {
					variable.ZapLog.Info("已经发送成功ID:" + strconv.Itoa(int(w.WsClient.HomeId)) + string(returnCalledStr))
				}
			}
		}

		return w, true
	} else {
		return nil, false
	}
}

// OnMessage 处理业务消息
func (w *Ws) OnMessage(context *gin.Context) {
	go w.WsClient.ReadPump(func(messageType int, receivedData []byte) {
		//参数说明
		//messageType 消息类型，1=文本
		//receivedData 服务器接收到客户端（例如js客户端）发来的的数据，[]byte 格式
		var receivedJson ReceivedWsClientMsg

		if err := json.Unmarshal(receivedData, &receivedJson); err == nil {
			switch receivedJson.Type {
			case 2:
				if w.CheckOnline(receivedJson.UserId) {

				}
			case 1: // 视频通话
				// 回复客户端已经收到消息
				// 被呼叫方在线
				if w.CheckOnline(receivedJson.UserId) {
					// 当前用户传来 code = 203 即挂断，向挂断方发送消息
					if receivedJson.Code == 203 {
						returnRefusedCall := ReturnClientMsg{}
						returnRefusedCall.Code = 203
						returnRefusedCall.Msg = "success"
						returnRefusedCall.Data.Type = 1
						returnRefusedCall.Data.UserId = w.WsClient.HomeId
						returnRefuseCallStr, _ := json.Marshal(returnRefusedCall)
						home_user.CreateHomeModelFactory("").UpdateIsCall(int(receivedJson.UserId), int(w.WsClient.HomeId), 0)
						// 向被挂断方定向发送信息
						w.SendMsgToClient(receivedJson.UserId, w.WsClient.UserType, string(returnRefuseCallStr))
					} else {
						// 呼叫方：当前用户 w.WsClient.HomeId
						// 被呼叫方：receivedJson.UserId
						// 给被呼叫方发送呼叫方信息以及房间号
						returnCalled := ReturnClientMsg{}
						returnCalled.Code = 200
						returnCalled.Msg = "success"
						returnCalled.Data.Type = 1
						returnCalled.Data.UserId = w.WsClient.HomeId
						roomData := room.CreateRoomModelFactory("").InsertData(int(w.WsClient.HomeId))
						w.WsClient.RoomId = roomData.Id
						returnCalled.Data.RoomId = roomData.Id

						// 查询被呼叫方是否有呼叫方好友
						friendData := friend_user.CreateFriendUserModelFactory("").GetByUserIdAndFriendId(int(receivedJson.UserId), int(w.WsClient.HomeId))
						userInfo := home_user.CreateHomeModelFactory("").GetHomeUser(w.WsClient.HomeId)
						returnCalled.Data.DeviceType = userInfo.DeviceType
						returnCalled.Data.UserIdCallerAvatar = userInfo.Avatar
						returnCalled.Data.HandFree = friendData.HandsFree
						if friendData.NickName != "" {
							returnCalled.Data.UserIdCallerTitle = friendData.NickName
						} else {
							returnCalled.Data.UserIdCallerTitle = userInfo.Title
						}
						returnCalledStr, _ := json.Marshal(returnCalled)

						// 这里处理一下呼叫状态，呼叫中的人，不能再被呼叫
						home_user.CreateHomeModelFactory("").UpdateIsCallOne(int(w.WsClient.HomeId), int(receivedJson.UserId))
						home_user.CreateHomeModelFactory("").UpdateIsCallOne(int(receivedJson.UserId), int(w.WsClient.HomeId))

						// 小程序要延迟三秒发送消息，因为 apk 进入程序有延迟
						calledUserInfo := home_user.CreateHomeModelFactory("").GetHomeUser(receivedJson.UserId)
						if w.WsClient.UserType == "web" && calledUserInfo.DeviceType == 2 {
							// 给呼叫方(当前用户)发送：给被呼叫方发送请求成功 202
							if err = w.WsClient.SendMessage(messageType, w.SendSuccess(202, int(receivedJson.UserId), "请求发送成功", roomData.Id, calledUserInfo.DeviceType)); err != nil {
								variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
							} else {
								variable.ZapLog.Info("已经发送成功ID:" + strconv.Itoa(int(w.WsClient.HomeId)) + w.SendSuccess(202, int(receivedJson.UserId), "请求发送成功", roomData.Id, calledUserInfo.DeviceType))
							}
							time.Sleep(4 * time.Second)
							// 给被呼叫方发送呼叫方信息以及房间号
							w.SendMsgToClient(receivedJson.UserId, w.WsClient.UserType, string(returnCalledStr))
						} else {
							w.SendMsgToClient(receivedJson.UserId, w.WsClient.UserType, string(returnCalledStr))
							if err = w.WsClient.SendMessage(messageType, w.SendSuccess(202, int(receivedJson.UserId), "请求发送成功", roomData.Id, calledUserInfo.DeviceType)); err != nil {
								variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
							} else {
								variable.ZapLog.Info("已经发送成功ID:" + strconv.Itoa(int(w.WsClient.HomeId)) + w.SendSuccess(202, int(receivedJson.UserId), "请求发送成功", roomData.Id, calledUserInfo.DeviceType))
							}
						}
					}
				} else { // 被呼叫方不在线
					// 给被呼叫方发送通话记录
					//calledFriendData := friend_user.CreateFriendUserModelFactory("").GetByUserIdAndFriendId(int(receivedJson.UserId), int(w.WsClient.HomeId))
					//callInfo := home_user.CreateHomeModelFactory("").GetHomeUser(w.WsClient.HomeId)
					userInfo := home_user.CreateHomeModelFactory("").GetHomeUser(receivedJson.UserId)

					// 如果是小程序，不论在不在线，都会请求发送成功，然后发布订阅消息
					if userInfo.DeviceType == 2 {

					} else {
						if err = w.WsClient.SendMessage(messageType, w.GetFail(201, "用户不在线")); err != nil {
							variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
						} else {
							variable.ZapLog.Info("已经发送成功ID:" + strconv.Itoa(int(w.WsClient.HomeId)) + w.GetFail(201, "用户不在线"))
						}
					}
				}
			case 35:
			default:
				variable.ZapLog.Info("收到ws客户端IP("+w.WsClient.ClientIp+")的消息", zap.String("msg", string(receivedData)))

			}
		}

		//tempMsg := "服务器已经收到了你的消息==>" + string(receivedData)
		//// 回复客户端已经收到消息;
		//if err := w.WsClient.SendMessage(messageType, tempMsg); err != nil {
		//	variable.ZapLog.Error("消息发送出现错误", zap.Error(err))
		//}

	}, w.OnError, w.OnClose)
}

// OnError 客户端与服务端在消息交互过程中发生错误回调函数
func (w *Ws) OnError(err error) {
	w.WsClient.State = 0 // 发生错误，状态设置为0, 心跳检测协程则自动退出
	variable.ZapLog.Error("远端掉线、卡死、刷新浏览器等会触发该错误:", zap.Error(err))
	//fmt.Printf("远端掉线、卡死、刷新浏览器等会触发该错误: %v\n", err.Error())
}

// OnClose 客户端关闭回调，发生onError回调以后会继续回调该函数
func (w *Ws) OnClose() {
	// 触发 onClose 时，需要把在播状态清除
	if w.WsClient.UserType != "web" {
		// 更新通话用户双方的状态
		homeId := w.WsClient.HomeId
		homeInfo := home_user.CreateHomeModelFactory("").GetHomeUser(homeId)
		fkFriendId := homeInfo.IsCall
		if fkFriendId != 0 {
			home_user.CreateHomeModelFactory("").UpdateIsCall(int(homeId), fkFriendId, 0)
		}

		// 更新通话记录的通话状态
		callLogData := call_log.CreateCallLogModelFactory("").GetByFkUserId(homeId)
		if callLogData.Id != 0 {
			callLogData.IsCall = 0
			// 修改当前用户的通话记录状态
			call_log.CreateCallLogModelFactory("").UpdateData(&callLogData)
			// 修改对方的通话记录状态
			call_log.CreateCallLogModelFactory("").UpdateIsCall(callLogData.FkFriendId, callLogData.FkUserId)
		}
	}
	variable.ZapLog.Info("用户下线；ID:" + strconv.Itoa(int(w.WsClient.HomeId)) + "类型：" + w.WsClient.UserType)
	w.WsClient.State = 0
	w.WsClient.Hub.UnRegister <- w.WsClient // 向hub管道投递一条注销消息，由hub中心负责关闭连接、删除在线数据
}

// GetOnlineClients  获取在线的全部客户端
func (w *Ws) GetOnlineClients() {

	fmt.Printf("在线客户端数量：%d\n", len(w.WsClient.Hub.Clients))
}

// BroadcastMsg  (每一个客户端都有能力)向全部在线客户端广播消息
func (w *Ws) BroadcastMsg(sendMsg string) {
	for onlineClient := range w.WsClient.Hub.Clients {

		//获取每一个在线的客户端，向远端发送消息
		if err := onlineClient.SendMessage(websocket.TextMessage, sendMsg); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
		}
	}
}

// CheckOnline 判断是否在线
func (w *Ws) CheckOnline(clientUserId int64) bool {
	for onlineClient := range w.WsClient.Hub.Clients {
		if clientUserId == onlineClient.ClientMoreParams.HomeId {
			return true
		}
	}
	return false
}

// SendMsgToClient 定向发送消息
func (w *Ws) SendMsgToClient(clientUserId int64, userType, sendMsg string) {
	for onlineClient := range w.WsClient.Hub.Clients {
		if onlineClient.ClientMoreParams.HomeId == clientUserId {
			// 获取每一个在线的客户端，向远端发送消息
			if userType == "web" || userType == "mobile" {
				if onlineClient.ClientMoreParams.UserType == "apk" || onlineClient.ClientMoreParams.UserType == "mobile" {
					if err := onlineClient.SendMessage(websocket.TextMessage, sendMsg); err != nil {
						variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
					} else {
						variable.ZapLog.Info("已经发送成功;对方ID:" + strconv.Itoa(int(clientUserId)) + sendMsg)
					}
					break
				}
			} else {
				if err := onlineClient.SendMessage(websocket.TextMessage, sendMsg); err != nil {
					variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
				} else {
					variable.ZapLog.Info("已经发送成功;对方ID:" + strconv.Itoa(int(clientUserId)) + sendMsg + "对方设备类型：" + onlineClient.ClientMoreParams.UserType)
				}
				continue
			}
		}
	}
}

func (w *Ws) GetFail(code int, msg string) string {
	returnData := ReturnFiled{}
	returnData.Code = int64(code)
	returnData.Msg = msg
	returnDataStr, _ := json.Marshal(returnData)
	return string(returnDataStr)
}

func (w *Ws) SendSuccess(code, userId int, msg string, roomId int64, deviceType int) string {
	returnData := ReturnClientMsg{}
	returnData.Code = int64(code)
	returnData.Msg = msg
	returnData.Data.Type = 1
	returnData.Data.UserId = int64(userId)
	returnData.Data.RoomId = roomId
	returnData.Data.DeviceType = deviceType
	returnDataStr, _ := json.Marshal(returnData)
	return string(returnDataStr)
}
