package job

import (
	"encoding/json"
	"fmt"
	"pusher/app/mp"
)

// Template struct 服务号模板消息推送
type Template struct {
	OpenID              string
	ReqID               string
	TemplateID          string
	PageURL             string
	MiniprogramAppID    string
	MiniprogramPagePath string
	ClientMSGID         string
	Key1                string
	Value1              string
	Color1              string
	Key2                string
	Value2              string
	Color2              string
	Key3                string
	Value3              string
	Color3              string
	Key4                string
	Value4              string
	Color4              string
	Key5                string
	Value5              string
	Color5              string
	AccessToken         string
}

// Run Player.Run 服务号模板消息推送最多5个参数
func (msg *Template) Run() {

	// log.Println(`Message Run11`)

	data := make(map[string]interface{})
	params := make(map[string]interface{})
	data["touser"] = msg.OpenID
	data["template_id"] = msg.TemplateID
	if msg.PageURL != "" {
		data["url"] = msg.PageURL
	}
	if msg.MiniprogramAppID != "" && msg.MiniprogramPagePath != "" {
		data["miniprogram"] = map[string]string{"appid": msg.MiniprogramAppID, "pagepath": msg.MiniprogramPagePath}
	}
	if msg.ClientMSGID != "" {
		data["client_msg_id"] = msg.ClientMSGID
	}
	params[msg.Key1] = map[string]string{"value": msg.Value1, "color": msg.Color1} //参数1
	if msg.Key2 != "" && msg.Value2 != "" {                                        // 参数2
		params[msg.Key2] = map[string]string{"value": msg.Value2, "color": msg.Color2}
	}
	if msg.Key3 != "" && msg.Value3 != "" { // 参数3
		params[msg.Key3] = map[string]string{"value": msg.Value3, "color": msg.Color3}
	}
	if msg.Key4 != "" && msg.Value4 != "" { // 参数4
		params[msg.Key4] = map[string]string{"value": msg.Value4, "color": msg.Color4}
	}
	if msg.Key5 != "" && msg.Value5 != "" { // 参数5
		params[msg.Key5] = map[string]string{"value": msg.Value5, "color": msg.Color5}
	}

	data["data"] = params

	datajson, _ := json.Marshal(data)

	urlStr := `https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=` + msg.AccessToken

	m := mp.CloudHttpPost(urlStr, datajson)

	if errmsg, ok := m["errmsg"]; ok {
		if errmsg == "ok" { // 推送失败
			fmt.Println(fmt.Sprintf("stok [%s]", msg.OpenID))
		} else {
			fmt.Println(fmt.Sprintf("stno [%s] [%s] [%s]", msg.OpenID, msg.TemplateID, errmsg))
		}
	} else { //推送异常
		fmt.Println(fmt.Sprintf("stff [%s] [%s]", msg.OpenID, msg.TemplateID))
		fmt.Println(m) // 这里是有问题的
	}

}
