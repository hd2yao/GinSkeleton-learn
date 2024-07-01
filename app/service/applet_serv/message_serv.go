package applet_serv

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/qifengzhang007/goCurl"
	"go.uber.org/zap"

	"goskeleton/app/global/variable"
)

func GetMessageServ() *MessageServ {
	return &MessageServ{}
}

type MessageServ struct{}

type TokenResponse struct {
	Errcode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
}

// GetToken 获取 Token
func (m *MessageServ) GetToken() (Token string, err error) {
	url := variable.ConfigYml.GetString("WeiXin.TokenUrl")

	cli := goCurl.CreateHttpClient(goCurl.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json; charset=utf-8",
		},
		Timeout: 10,
	})

	resp, err := cli.Get(url, goCurl.Options{
		FormParams: map[string]interface{}{
			"grant_type": "client_credential",
			"appid":      variable.ConfigYml.GetString("WeiXin.Appid"),
			"secret":     variable.ConfigYml.GetString("WeiXin.Secret"),
		},
		SetResCharset: "utf-8",
	})

	if err != nil {
		variable.ZapLog.Error("向腾讯获取token出错：", zap.Error(err))
		return
	} else {
		var res string
		var received TokenResponse
		res, err = resp.GetContents()
		if err == nil {
			if err = json.Unmarshal([]byte(res), &received); err == nil {
				if received.Errcode == 0 {
					Token = received.AccessToken
					return Token, nil
				} else {
					variable.ZapLog.Sugar().Infof("向腾讯获取token出错：%+v", received)
					return
				}
			} else {
				variable.ZapLog.Error("腾讯返回数据反序列化出错", zap.Error(err))
				return
			}
		}
	}
	return
}

type MsgRequestDataType struct {
	Value string `json:"value"`
}

type MsgRequest struct {
	ToUser           string                        `json:"touser"`
	TemplateID       string                        `json:"template_id"`
	Page             string                        `json:"page"`
	Lang             string                        `json:"lang"`
	MiniprogramState string                        `json:"miniprogram_state"`
	Data             map[string]MsgRequestDataType `json:"data"`
}

type MsgResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (m *MessageServ) SendMessage(openid, title, friendTitle string) error {
	//先到数据库查询token，因为腾讯token过期时间是两小时，所以做一下存储
	token := ""
	// 数据库相关操作 这里暂时注释掉
	//tokenData := applet_token.CreateAppletTokenModelFactory("").GetToken()
	//if tokenData.Id != 0 {
	//	token = tokenData.Token
	//} else {
	//	tokenGet, err := m.GetToken()
	//	if err != nil {
	//		return err
	//	} else {
	//		token = tokenGet
	//		//存入数据库一份
	//		insertToken := applet_token.AppletTokenModel{}
	//		insertToken.Token = tokenGet
	//		applet_token.CreateAppletTokenModelFactory("").InsertData(&insertToken)
	//	}
	//}

	url := variable.ConfigYml.GetString("WeiXin.MsgUrl")
	url = url + "?access_token=" + token
	cli := goCurl.CreateHttpClient(goCurl.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json; charset=utf-8",
		},
		SetResCharset: "utf-8",
		Timeout:       10,
	})
	data := make(map[string]MsgRequestDataType)
	data["thing1"] = MsgRequestDataType{Value: friendTitle}
	data["thing2"] = MsgRequestDataType{Value: title}
	data["time3"] = MsgRequestDataType{Value: time.Now().Format(variable.DateFormat)}

	var sendData = MsgRequest{
		ToUser:           openid,
		TemplateID:       variable.ConfigYml.GetString("WeiXin.TemplateId"),
		Page:             "pages/index",
		Lang:             "zh_CN",
		MiniprogramState: variable.ConfigYml.GetString("WeiXin.State"),
		Data:             data,
	}
	resp, err := cli.Post(url, goCurl.Options{
		JSON: sendData,
	})

	if err != nil {
		variable.ZapLog.Error("发布消息出错：", zap.Error(err))
		return err
	} else {
		var res string
		var received MsgResponse
		res, err = resp.GetContents()
		if err == nil {
			if err = json.Unmarshal([]byte(res), &received); err == nil {
				if received.Errcode == 0 {
					return nil
				} else {
					variable.ZapLog.Sugar().Infof("消息推送出错：%+v", received)
					return errors.New("发布消息出错")
				}
			} else {
				variable.ZapLog.Error("腾讯返回数据反序列化出错", zap.Error(err))
			}
		}
	}
	return nil
}
