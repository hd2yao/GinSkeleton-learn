package home_user

import "goskeleton/app/model"

func CreateHomeModelFactory(sqlType string) *HomeUserModel {
	return &HomeUserModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type HomeUserModel struct {
	model.BaseModel
	Title            string `json:"title"`
	Mac              string `json:"mac"`
	Addr             string `json:"addr"`
	Openid           string `json:"openid"`
	HouseId          string `json:"house_id"`
	Pass             string `json:"pass"`
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	Avatar           string `json:"avatar"`
	DeviceType       int    `json:"device_type"`
	Tag              int    `json:"tag"`
	IsCall           int    `json:"is_call"`
	FkProvinceCityId int    `json:"fk_province_city_id"`
	FkPackageId      int    `json:"fk_package_id"`
	Remark           string `json:"remark"`
}

func (h *HomeUserModel) TableName() string {
	return "tb_home"
}

// UpdateIsCall 拨打中的状态变更
// is_call = 0 表示当前不在通话中；is_call 不为 0 时，is_call == fk_friend_id 表示正在通话的用户 id
func (h *HomeUserModel) UpdateIsCall(userId1, userId2, isCall int) {
	h.Model(h).Where("id = ? or id = ?", userId1, userId2).Update("is_call", isCall)
}

// UpdateIsCallToZero 更新用户通话状态为不在通话中
func (h *HomeUserModel) UpdateIsCallToZero(userId int) {
	h.Model(h).Where("id = ?", userId).Update("is_call", 0)
}

// GetHomeUser 查询指定用户信息(网关或小程序用户)
func (h *HomeUserModel) GetHomeUser(userId int64) (data HomeUserModel) {
	data = HomeUserModel{}
	h.Model(h).Where("id = ?", userId).Scan(&data)
	return data
}

// GetCallId 查询用户正在通话的用户
func (h *HomeUserModel) GetCallId(userId int) int {
	data := HomeUserModel{}
	h.Model(h).Where("id = ?", userId).First(&data)
	return data.IsCall
}
