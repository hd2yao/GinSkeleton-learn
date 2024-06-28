package msg

import "goskeleton/app/model"

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
