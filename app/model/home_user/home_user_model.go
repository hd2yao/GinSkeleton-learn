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
func (h *HomeUserModel) UpdateIsCall(userId1, userId2, isCall int) {
    h.Model(h).Where("id = ? or id = ?", userId1, userId2).Update("is_call", isCall)
}
