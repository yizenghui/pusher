package mp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

// httpClient 默认http.Client
var httpClient *http.Client

func init() {
	client := *http.DefaultClient
	client.Timeout = time.Second * 5
	httpClient = &client
}

// OpenIDData 开放数据 openID
type OpenIDData struct {
	ErrCode    int64  `json:"errcode"`
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

// GetOpenID 获取微信小程序上报的openid 此ID暂不加密处理
func GetOpenID(code string) (OpenIDData, error) {
	//
	type Ret struct {
		ErrCode    int64  `json:"errcode"`
		ErrMSG     string `json:"errmsg"`
		SessionKey string `json:"session_key"`
		ExpiresIn  int64  `json:"expires_in"`
		OpenID     string `json:"openid"`
	}
	var ret Ret

	AppID := os.Getenv("MINAPP_APPID")
	AppSecret := os.Getenv("MINAPP_APPSECRET")

	url := fmt.Sprintf(`https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code`,
		AppID,
		AppSecret,
		code,
	)

	HTTPGetJSON(url, &ret)
	var err error

	if ret.ErrCode != 0 {
		err = errors.New(fmt.Sprint(ret.ErrCode))
	}

	// log.Println("Ret", ret)
	return OpenIDData{ret.ErrCode, ret.OpenID, ret.SessionKey}, err
}

// HTTPGetJSON 通过传入url和结构，提取出页面中的值
func HTTPGetJSON(url string, response interface{}) error {
	httpResp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	return DecodeJSONHttpResponse(httpResp.Body, response)
}

// HTTPPostJSON  通过传入url和内容，提交内容后，提取出页面中的值
func HTTPPostJSON(url string, body []byte, response interface{}) error {
	httpResp, err := httpClient.Post(url, "application/json; charset=utf-8", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	return DecodeJSONHttpResponse(httpResp.Body, response)
}

// DecodeJSONHttpResponse 解决json
func DecodeJSONHttpResponse(r io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// CloudHttpGet 类php常用get函数
func CloudHttpGet(url string) map[string]interface{} {

	resp, _ := httpClient.Get(url)

	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	m, _ := gjson.ParseBytes([]byte(body)).Value().(map[string]interface{})

	return m
}

// CloudHttpPost 类php常用post函数
func CloudHttpPost(url string, params []byte) map[string]interface{} {

	resp, _ := httpClient.Post(url, "application/json; charset=utf-8", bytes.NewReader(params))

	// log.Println(`CloudHttpPost Run`, resp)
	body, _ := ioutil.ReadAll(resp.Body)

	// log.Println(`CloudHttpPost Run`, string(body))
	defer resp.Body.Close()

	m, _ := gjson.ParseBytes([]byte(body)).Value().(map[string]interface{})

	return m
}
