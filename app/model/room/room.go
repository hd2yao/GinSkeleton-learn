package room

import "goskeleton/app/model"

func CreateRoomModelFactory(sqlType string) *RoomModel {
	return &RoomModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type RoomModel struct {
	model.BaseModel
	Status   int    `json:"status"`
	FkUserId int    `json:"fk_user_id"`
	Remark   string `json:"remark"`
}

func (r *RoomModel) TableName() string {
	return "tb_room"
}

func (r *RoomModel) InsertData(userId int) RoomModel {
	data := RoomModel{}
	data.Status = 1
	data.FkUserId = userId
	r.Create(&data)
	return data
}

// GetRoomId 获取roomId
func (r *RoomModel) GetRoomId(userId int) int64 {
	data := RoomModel{}
	r.Model(r).Where("fk_user_id = ?", userId).Order("id desc").First(&data)
	return data.Id
}
