package video_data

type ActiveVali struct {
    FkUserId   int `form:"fk_user_id" json:"fk_user_id" binding:"required"`
    FkFriendId int `form:"fk_friend_id" json:"fk_friend_id" binding:"required"`
    IsCall     int `form:"is_call" json:"is_call"`
    Types      int `form:"types" json:"types" binding:"required"`
}
