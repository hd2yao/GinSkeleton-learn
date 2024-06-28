package msg

import "goskeleton/app/model"

func CreateMsgContentModelFactory(sqlType string) *MsgContentModel {
    return &MsgContentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type MsgContentModel struct {
    model.BaseModel
    Title            string `json:"title"`
    StartTime        string `json:"start_time"`
    EndTime          string `json:"end_time"`
    Content          string `json:"content"`
    FkProvinceCityId int    `json:"fk_province_city_id"`
    Remark           string `json:"remark"`
}

func (m *MsgContentModel) TableName() string {
    return "tb_msg_content"
}

func (m *MsgContentModel) GetContent(cityId int) []MsgContentModel {
    data := []MsgContentModel{}
    m.Model(m).Where("NOW()>start_time").
        Where("NOW()<end_time").
        Where("FIND_IN_SET(fk_province_city_id,(SELECT path_info FROM tb_province_city WHERE id = ?))", cityId).Find(&data)
    return data
}
