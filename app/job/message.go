package job

import (
	"encoding/json"
	"fmt"

	"github.com/yizenghui/pusher/app/mp"
)

// Message struct 小程序订阅消息推送
type Message struct {
	OpenID      string
	ReqID       string
	TemplateID  string
	Page        string
	Key1        string
	Value1      string
	Key2        string
	Value2      string
	Key3        string
	Value3      string
	Key4        string
	Value4      string
	Key5        string
	Value5      string
	AccessToken string
}

// Run Player.Run 小程序模板消息推送最多5个参数
func (msg *Message) Run() {

	// log.Println(`Message Run11`)

	data := make(map[string]interface{})
	params := make(map[string]interface{})
	data["touser"] = msg.OpenID
	data["template_id"] = msg.TemplateID
	data["page"] = msg.Page
	data["miniprogram_state"] = "formal"
	params[msg.Key1] = map[string]string{"value": msg.Value1}
	if msg.Key2 != "" && msg.Value2 != "" { // 参数2
		params[msg.Key2] = map[string]string{"value": msg.Value2}
	}
	if msg.Key3 != "" && msg.Value3 != "" { // 参数3
		params[msg.Key3] = map[string]string{"value": msg.Value3}
	}
	if msg.Key4 != "" && msg.Value4 != "" { // 参数4
		params[msg.Key4] = map[string]string{"value": msg.Value4}
	}
	if msg.Key5 != "" && msg.Value5 != "" { // 参数5
		params[msg.Key5] = map[string]string{"value": msg.Value5}
	}

	data["data"] = params

	datajson, _ := json.Marshal(data)

	urlStr := `https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=` + msg.AccessToken

	m := mp.CloudHttpPost(urlStr, datajson)

	if errmsg, ok := m["errmsg"]; ok {
		if errmsg == "ok" { // 推送失败
			fmt.Println(fmt.Sprintf("ssok [%s]", msg.OpenID))
		} else {
			fmt.Println(fmt.Sprintf("ssno [%s] [%s] [%s]", msg.OpenID, msg.TemplateID, errmsg))
		}
	} else { //推送异常
		fmt.Println(fmt.Sprintf("ssff [%s] [%s]", msg.OpenID, msg.TemplateID))
		fmt.Println(m) // 这里是有问题的
	}

}
