package call_log

import (
    "go.uber.org/zap"
    "goskeleton/app/global/variable"
    "goskeleton/app/model"
)

func CreateCallLogModelFactory(sqlType string) *CallLogModel {
    return &CallLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type CallLogModel struct {
    model.BaseModel
    FkUserId     int    `json:"fk_user_id"`
    FkCallTypeId int    `json:"fk_call_type_id"`
    FkFriendId   int    `json:"fk_friend_id"`
    IsCall       int    `json:"is_call"`
    CallStatus   int    `json:"call_status"`
    CallDuration int    `json:"call_duration"`
    ConnectTime  int    `json:"connect_time"`
    RingOffTime  int    `json:"ring_off_time"`
    IsRead       int    `json:"is_read"`
    Remark       string `json:"remark"`
}

func (c *CallLogModel) TableName() string {
    return "tb_call_log"
}

func (c *CallLogModel) InsertData(formData *CallLogModel) bool {
    if res := c.Create(formData); res.Error != nil {
        variable.ZapLog.Error("CallLogModel 数据新增出错", zap.Error(res.Error))
        return false
    }
    return true
}
