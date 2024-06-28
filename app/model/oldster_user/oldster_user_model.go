package oldster_user

import "goskeleton/app/model"

func CreateOldsterUserModelFactory(sqlType string) *OldsterUserModel {
    return &OldsterUserModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type OldsterUserModel struct {
    model.BaseModel
    Name             string `json:"name"`
    CardId           string `json:"card_id"`
    CardIdMD5        string `json:"card_id_md5"`
    Avatar           string `json:"avatar"`
    TagIds           string `json:"tag_ids"`
    FkProvinceCityId int    `json:"fk_province_city_id"`
    Phone            string `json:"phone"`
    Addr             string `json:"addr"`
    FkHomeIds        string `json:"fk_home_ids"`
    Remark           string `json:"remark"`
}

func (o *OldsterUserModel) TableName() string {
    return "tb_oldster_user"
}

// GetById 根据 id 查询老人信息
func (o *OldsterUserModel) GetById(id int64) (data OldsterUserModel) {
    o.Model(o).Where("id = ?", id).First(&data)
    return
}
