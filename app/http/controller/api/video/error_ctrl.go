package video

import (
    "github.com/gin-gonic/gin"

    "goskeleton/app/global/consts"
    "goskeleton/app/model/call_log"
    "goskeleton/app/model/home_user"
    "goskeleton/app/utils/response"
)

type ErrorType struct{}

type ErrorVali struct {
    HomeId int `form:"home_id" json:"home_id" binding:"required"`
}

func (e ErrorType) HandleError(context *gin.Context) {
    form := ErrorVali{}
    if err := context.ShouldBind(&form); err != nil {
        response.ValidatorError(context, err)
        return
    }

    homeUserData := home_user.CreateHomeModelFactory("").GetHomeUser(int64(form.HomeId))
    fkFriendId := homeUserData.IsCall

    // 更新用户双方的通话状态
    if fkFriendId != 0 {
        // 更新自己的通话状态为 0
        home_user.CreateHomeModelFactory("").UpdateIsCallToZero(form.HomeId)
        // 获取好友信息
        friendData := home_user.CreateHomeModelFactory("").GetHomeUser(int64(fkFriendId))
        if friendData.IsCall == form.HomeId {
            // 也可以不需要获取好友信息后再修改，这样做多了一层保险
            home_user.CreateHomeModelFactory("").UpdateIsCallToZero(fkFriendId)
        }
    }

    // 更新通话记录的状态
    latestLog := call_log.CreateCallLogModelFactory("").FindIsCallById(form.HomeId)
    if latestLog.Id != 0 {
        latestLog.IsCall = 0
        // 更新当前用户的通话记录状态
        call_log.CreateCallLogModelFactory("").UpdateData(&latestLog)
        // 更新对方的通话记录状态
        call_log.CreateCallLogModelFactory("").UpdateIsCall(latestLog.FkFriendId, latestLog.FkUserId)
    }
    response.Success(context, consts.CurdStatusOkMsg, "")
}
