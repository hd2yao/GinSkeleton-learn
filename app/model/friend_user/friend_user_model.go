package friend_user

import "goskeleton/app/model"

func CreateFriendUserModelFactory(sqlType string) *FriendUserModel {
    return &FriendUserModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type FriendUserModel struct {
    model.BaseModel
    FkUserId   int    `json:"fk_user_id"`
    FkFriendId int    `json:"fk_friend_id"`
    NickName   string `json:"nick_name"`
    Top        int    `json:"top"`
    HandsFree  int    `json:"hands_free"`
    Remark     string `json:"remark"`
}

func (f *FriendUserModel) TableName() string {
    return "tb_friend"
}

// GetByUserIdAndFriendId 获取用户指定好友信息
func (f *FriendUserModel) GetByUserIdAndFriendId(fkUserId, fkFriendId int) FriendUserModel {
    data := FriendUserModel{}
    f.Model(f).Where("fk_user_id = ?", fkUserId).Where("fk_friend_id = ?", fkFriendId).First(&data)
    return data
}
