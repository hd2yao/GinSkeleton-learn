package websocket

import (
    "goskeleton/app/model/msg"
    "goskeleton/app/model/oldster_user"
    "goskeleton/app/utils/websocket/core"
)

// 这里可以定义业务相关的逻辑供 ws.go 文件调用

func HandleMsg(cli *core.Client, isOpen bool) {
    //只给web发送信息
    if cli.UserType != "web" {
        return
    }

    // 循环在线的用户，去查询是否有需要推送内容的用户
    // 1. 先查看用户需要推送的信息
    msgContents := msg.CreateMsgContentModelFactory("").GetContent(int(cli.CityId))

    // 查看里面的消息是否已经推送过，推送过就不需要再推，去推送记录表查看
    // tag 用于标识是否需要向前端推送消息
    tag := false
    userInfo := oldster_user.CreateOldsterUserModelFactory("").GetById(cli.UserId)
    for _, content := range msgContents {
        recordData := msg.CreateMsgRecordModelFactory("").GetByUserIdAndContentId(cli.UserId, content.Id)
        if recordData.Id == 0 {
            // 如果没有记录，就增加一条记录
            recordData.FkUserId = int(cli.UserId)
            recordData.FkUserName = userInfo.Name
            recordData.Content = content.Content
            recordData.Title = content.Title
            recordData.FkMsgContentId = int(content.Id)
            msg.CreateMsgRecordModelFactory("").InsertData(&recordData)
            tag = true
        } else {
            // 推送内容有修改，要重新推送
            if recordData.Content != content.Content {
                recordData.Content = content.Content
                msg.CreateMsgRecordModelFactory("").InsertData(&recordData)
                tag = true
            }
        }
    }

}
