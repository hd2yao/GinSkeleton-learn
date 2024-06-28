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

// GetContent 获取当前区域内需要推送的消息内容
func (m *MsgContentModel) GetContent(cityId int) []MsgContentModel {
	data := []MsgContentModel{}
	m.Model(m).Where("NOW()>start_time").
		Where("NOW()<end_time").
		Where("FIND_IN_SET(fk_province_city_id,(SELECT path_info FROM tb_province_city WHERE id = ?))", cityId).Find(&data)
	return data
}

// BirthdayRecord 查询当前用户是否有今年的生日推送记录
func (m *MsgRecordModel) BirthdayRecord(userId int64) bool {
	data := MsgRecordModel{}
	m.Model(m).Where("fk_msg_content_id = 0").Where("YEAR(CURDATE()) = DATE_FORMAT(created_at,'%Y')").Where("fk_user_id = ?", userId).First(&data)
	if data.Id == 0 {
		return false
	}
	return true
}
