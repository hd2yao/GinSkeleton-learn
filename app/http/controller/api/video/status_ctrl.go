package video

import (
    "github.com/gin-gonic/gin"

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
