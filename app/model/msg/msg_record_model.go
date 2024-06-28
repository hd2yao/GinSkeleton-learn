package msg

import (
    "go.uber.org/zap"

    "goskeleton/app/global/variable"
    "goskeleton/app/model"
)

func CreateMsgRecordModelFactory(sqlType string) *MsgRecordModel {
    return &MsgRecordModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type MsgRecordModel struct {
    model.BaseModel
    Title          string `json:"title"`
    Content        string `json:"content"`
    FkUserId       int    `json:"fk_user_id"`
    FkUserName     string `json:"fk_user_name"`
    FkMsgContentId int    `json:"fk_msg_content_id"`
    Remark         string `json:"remark"`
}

func (m *MsgRecordModel) TableName() string {
    return "tb_msg_record"
}

func (m *MsgRecordModel) GetByUserIdAndContentId(userId, contentId int64) (data MsgRecordModel) {
    m.Model(m).Where("fk_user_id = ?", userId).Where("fk_msg_content_id = ?", contentId).Find(&data)
    return
}

func (m *MsgRecordModel) InsertData(formData *MsgRecordModel) bool {
    if res := m.Create(formData); res.Error != nil {
        variable.ZapLog.Error("MsgMediaModel 数据新增出错", zap.Error(res.Error))
        return false
    }
    return true
}
