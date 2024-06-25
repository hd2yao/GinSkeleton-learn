package video

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/call_log"
	"goskeleton/app/model/home_user"
	"time"

	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/api/video/video_data"
	"goskeleton/app/service/call_status"
	"goskeleton/app/utils/response"
)

type VideoStatus struct{}

// 插入通话记录

func (v VideoStatus) HandleStatus(context *gin.Context) {
	form := video_data.ActiveVali{}
	if err := context.ShouldBind(&form); err != nil {
		response.ValidatorError(context, err)
		return
	}
	err := call_status.GetCallStatusService().Handle(form)
	if err == nil {
		response.Success(context, consts.CurdStatusOkMsg, "")
	}
}

// 通话中挂断

type RingOffVali struct {
	FkUserId     int `form:"fk_user_id" json:"fk_user_id" binding:"required"`
	FkFriendId   int `form:"fk_friend_id" json:"fk_friend_id" binding:"required"`
	Types        int `form:"types" json:"types" binding:"required"`
	IsCall       int `form:"is_call" json:"is_call"`
	CallDuration int `form:"call_duration" json:"call_duration" binding:"required"`
}

// 这个接口是轮询的，呼叫方和接听方都会调用，用于同步通话记录时长

func (v VideoStatus) RingOff(context *gin.Context) {
	form := RingOffVali{}
	if err := context.ShouldBind(&form); err != nil {
		response.ValidatorError(context, err)
		return
	}
	variable.ZapLog.Sugar().Infof("收到的数据：%+v", form)

	// 先查询是否有相关的通话记录
	data := call_log.CreateCallLogModelFactory("").FindById(form.FkUserId, form.FkFriendId, form.Types)
	if data.Id == 0 {
		response.Fail(context, consts.CurdSelectFailCode, "信息有误，没有相关通话记录", "")
		return
	}
	// 同步更新通话时长
	timestamp := time.Now().Unix()
	data.RingOffTime = int(timestamp)
	data.CallDuration = form.CallDuration
	data.IsCall = form.IsCall
	call_log.CreateCallLogModelFactory("").UpdateData(&data)

	// 在最后一次传入 isCall = 0，更新呼叫方的通话
	if data.IsCall == 0 {
		call_log.CreateCallLogModelFactory("").UpdateIsCall(form.FkFriendId, form.FkUserId)
		// 更改网关的通话状态
		home_user.CreateHomeModelFactory("").UpdateIsCall(form.FkUserId, form.FkFriendId, 0)
	}
	response.Success(context, consts.CurdStatusOkMsg, "")
}

// 查看是否在通话中（是否占线）

type IsCall struct {
	FkFriendId int `form:"fk_friend_id" json:"fk_friend_id" binding:"required"`
}

type IsCallReturn struct {
	IsCall bool `json:"is_call"`
}

func (v VideoStatus) IsCall(context *gin.Context) {
	form := IsCall{}
	if err := context.ShouldBind(&form); err != nil {
		response.ValidatorError(context, err)
		return
	}
	// 查询通话记录
	callLogData := call_log.CreateCallLogModelFactory("").FindIsCallById(form.FkFriendId)
	// 查询用户通话状态
	homeUserData := home_user.CreateHomeModelFactory("").GetHomeUser(int64(form.FkFriendId))
	returnData := IsCallReturn{}
	if callLogData.Id != 0 || homeUserData.IsCall != 0 {
		returnData.IsCall = true
	} else {
		returnData.IsCall = false
	}
	response.Success(context, consts.CurdStatusOkMsg, returnData)
}
