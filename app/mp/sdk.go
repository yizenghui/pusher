package mp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// MsgSecCheck 解决json
func MsgSecCheck(checktext string) int16 {

	var TokenServe *DefaultAccessTokenServer

	TokenServe = NewDefaultAccessTokenServer(`wx96830e80b331c267`, `e72aed927a8ccda32c2573705c7d098f`)
	token, err := TokenServe.Token()

	type CheckTextResult struct {
		ErrCode int64                    `json:"errcode"`
		ErrMsg  string                   `json:"errmsg"`
		TraceID string                   `json:"trace_id"`
		Result  int                      `json:"result"`
		Detail  []map[string]interface{} `json:"detail"`
	}

	urlStr := `https://api.weixin.qq.com/wxa/msg_sec_check?access_token=` + token

	data := `{"content":"特3456书 yuuo 莞6543李 zxcz 蒜7782法 fgnv 级
	完2347全 dfji 试3726测 asad 感3847知 qwez 到
"}` //习近平

	resp, err := http.Post(urlStr, "application/json", strings.NewReader(data))
	if err != nil {
		return -1
	}
	body, err := ioutil.ReadAll(resp.Body)

	// reader, err := charset.NewReader(resp.Body, strings.ToLower(resp.Header.Get("Content-Type")))
	defer resp.Body.Close()
	// bs, _ := ioutil.ReadAll(reader)

	var m CheckTextResult
	err = json.Unmarshal(body, &m)

	log.Println(m)
	return 0
}

// 详情看官方文档
type MediaCheckData struct {
	MediaURL  string `json:"media_url"`
	MediaType int32  `json:"media_type"` //  1:音频;2:图片
	Version   int32  `json:"version"`    //2
	Openid    string `json:"openid"`
	Scene     int32  `json:"scene"` // 2
}

// 媒体文件检查
func MediaSecCheck(medaiURL string, mediaType int32, openid string, token string) map[string]interface{} {

	urlStr := `https://api.weixin.qq.com/wxa/media_check_async?access_token=` + token

	data := &MediaCheckData{
		MediaURL:  medaiURL,
		MediaType: mediaType,
		Version:   2,
		Openid:    openid,
		Scene:     2,
	}

	datajson, _ := json.Marshal(data)

	return CloudHttpPost(urlStr, datajson)
}

// //WxAppSubmitPage 提交页面单个页面（望收录）
// func WxAppSubmitPage(checktext string) error {
// 	type Data struct {
// 		Content string `json:"content"`
// 	}
// 	//
// 	type CheckTextResult struct {
// 		ErrCode int64                    `json:"errcode"`
// 		ErrMsg  string                   `json:"errmsg"`
// 		TraceID string                   `json:"trace_id"`
// 		Result  int                      `json:"result"`
// 		Detail  []map[string]interface{} `json:"detail"`
// 	}

// 	var ret Ret

// 	var data = Data{}

// 	token, err2 := TokenServe.Token()
// 	if err2 != nil {
// 		return err2
// 	}
// 	// token = `27_EFpACLm1qpGcK8p_xEnZPnowJGKKEfWzy7500PLAR7Ek-8UaooSW-HTteSCfM2_r2f3zkKTcCgLFYvE094UNzXhZyv3KbZqAk_D8USQGFeYqklXrC6UVBIZfO0oAI2yB63nI0-cAsHjksNcAOPNjAEACDB`
// 	url := fmt.Sprintf(`https://api.weixin.qq.com/wxa/search/wxaapi_submitpages?access_token=%v`, token)

// 	b, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	HTTPPostJSON(url, b, &ret)

// 	// log.Println(`xx`, ret)
// 	if ret.ErrCode != 0 {
// 		// err = errors.New(string(ret.ErrMSG))
// 		err = errors.New(strconv.FormatInt(ret.ErrCode, 10))
// 	}

// 	return err
// }
