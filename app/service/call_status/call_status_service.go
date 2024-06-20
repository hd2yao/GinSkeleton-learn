package call_status

import (
    "strconv"
    "time"

    "goskeleton/app/http/controller/api/video/video_data"
    "goskeleton/app/model/call_log"
    "goskeleton/app/model/home_user"
    "goskeleton/app/service/websocket"
)

type CallStatusService struct{}

func GetCallStatusService() *CallStatusService {
    return &CallStatusService{}
}

func (c *CallStatusService) Handle(form video_data.ActiveVali) error {
    data := call_log.CallLogModel{}
    data.FkUserId = form.FkUserId
    data.FkFriendId = form.FkFriendId
    data.IsCall = 1

    switch form.Types {
    case 1:
        data.FkCallTypeId = 1
        data.CallStatus = 1
        data.IsRead = 1
        timestamp := time.Now().Unix()
        data.ConnectTime = int(timestamp)
    case 2:
        data.FkCallTypeId = 2
        data.CallStatus = 2
        data.IsRead = 1
        timestamp := time.Now().Unix()
        data.ConnectTime = int(timestamp)
    case 3:
        data.FkCallTypeId = 3
        timestamp := time.Now().Unix()
        data.RingOffTime = int(timestamp)
        data.IsCall = 0
        data.CallStatus = 1
        // 未接听的一方
        data2 := call_log.CallLogModel{}
        data2.FkUserId = form.FkFriendId
        data2.FkFriendId = form.FkUserId
        data2.FkCallTypeId = 3
        data2.RingOffTime = int(timestamp)
        data2.IsCall = 0
        data2.CallStatus = 2
        call_log.CreateCallLogModelFactory("").InsertData(&data2)
        // 通过 ws 向未接听一方推送通话记录
        userIdStr := strconv.Itoa(data2.FkUserId)
        go websocket.Pub(userIdStr)
    }
    // 修改网关双方的通话状态
    home_user.CreateHomeModelFactory("").UpdateIsCall(form.FkUserId, form.FkFriendId, 0)
    call_log.CreateCallLogModelFactory("").InsertData(&data)
    return nil
}
