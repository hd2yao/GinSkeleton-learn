package websocket

import (
    "encoding/json"

    "go.uber.org/zap"

    "goskeleton/app/global/variable"
    "goskeleton/app/model/msg"
    "goskeleton/app/model/oldster_user"
    "goskeleton/app/utils/websocket/core"
)

// 这里可以定义业务相关的逻辑供 ws.go 文件调用

// HandleMsg 推送内容业务（主动关爱）
func HandleMsg(cli *core.Client, isOpen bool) {
    // isOpen = true 表示 ws 连接时触发 onOpen 事件，会将全部的推送消息传给前端，用于首页滚动显示，一次连接只推送一次
    // isOpen = false 表示 ws 连接后实时监听，当有修改时，推送给前端

    // 只给web发送信息
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

    // 2. 查询用户是否今天生日
    isBirthday := oldster_user.CreateOldsterUserModelFactory("").GetIsBirthday(cli.UserId)
    if isBirthday {
        birthdayData := msg.MsgContentModel{}
        birthdayData.Title = "生日"
        birthdayData.Content = "今天是您的生日！祝您生日快乐, 健康幸福!"
        msgContents = append(msgContents, birthdayData)
        // 查看是否已经有推送过生日祝福
        hasBirthdayRecord := msg.CreateMsgRecordModelFactory("").BirthdayRecord(cli.UserId)
        if !hasBirthdayRecord {
            // 添加生日记录
            recordBirthdayData := msg.MsgRecordModel{}
            recordBirthdayData.FkUserId = int(cli.UserId)
            recordBirthdayData.FkUserName = userInfo.Name
            recordBirthdayData.Content = "今天是您的生日！祝您生日快乐, 健康幸福!"
            recordBirthdayData.Title = "生日"
            msg.CreateMsgRecordModelFactory("").InsertData(&recordBirthdayData)
            tag = true
        }
    }

    if cli.MsgLen != len(msgContents) {
        tag = true
        cli.MsgLen = len(msgContents)
    }

    // 如果缺少推送记录 或是 要推送全部信息，则需要给前端推送数据
    if tag == true || isOpen == true {
        sendMsg := MsgPush{}
        sendMsg.Code = 250
        sendMsg.Data = msgContents
        sendMsgMarshal, _ := json.Marshal(sendMsg)

        if err := cli.SendMessage(1, string(sendMsgMarshal)); err != nil {
            variable.ZapLog.Error("websocket定向发送消息出错", zap.Error(err))
        } else {
            variable.ZapLog.Info("websocket定向发送消息成功")
            tag = false
        }
    }
}
