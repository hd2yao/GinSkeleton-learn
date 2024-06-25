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

// InsertData 新增
func (c *CallLogModel) InsertData(formData *CallLogModel) bool {
	if res := c.Create(formData); res.Error != nil {
		variable.ZapLog.Error("CallLogModel 数据新增出错", zap.Error(res.Error))
		return false
	}
	return true
}

// UpdateData 更新
func (c *CallLogModel) UpdateData(formData *CallLogModel) bool {
	if res := c.Omit("created_at").Save(formData); res.Error != nil {
		variable.ZapLog.Error("CallLogModel 修改新增出错", zap.Error(res.Error))
		return false
	}
	return true
}

// FindById 根据通话双方 id 和通话类型查询正在通话中的通话记录
func (c *CallLogModel) FindById(fkUserId, fkFriendId, types int) (data CallLogModel) {
	data = CallLogModel{}
	c.Model(c).Where("fk_user_id = ?", fkUserId).
		Where("fk_friend_id = ?", fkFriendId).
		Where("fk_call_type_id = ?", types).
		Where("is_call = 1").
		Order("id desc").First(&data)
	return data
}

// FindIsCallById 查询指定用户是否正在通话中
// 如果用户正在通话中，则会有两条记录，只返回其中一条即可
func (c *CallLogModel) FindIsCallById(fkFriendId int) (data CallLogModel) {
	data = CallLogModel{}
	c.Model(c).
		Where("fk_user_id = ? or fk_friend_id = ?", fkFriendId, fkFriendId).
		Where("is_call = 1").
		Order("id desc").First(&data)
	return data
}

// UpdateIsCall 根据通话双方 id 将指定通话记录的状态从通话中更改为不在不在通话中
func (c *CallLogModel) UpdateIsCall(fkUserId, fkFriendId int) {
	c.Model(c).Where("fk_user_id = ?", fkUserId).
		Where("fk_friend_id = ?", fkFriendId).
		Where("is_call = 1").
		Update("is_call", 0)
}
