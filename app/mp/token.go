package mp

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"
)

var DefaultHttpClient *http.Client

func init() {
	client := *http.DefaultClient
	client.Timeout = time.Second * 5
	DefaultHttpClient = &client
}

type Error struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

const (
	ErrCodeOK                 = 0
	ErrCodeInvalidCredential  = 40001 // access_token 过期错误码
	ErrCodeAccessTokenExpired = 42001 // access_token 过期错误码(maybe!!!)
)

func (err *Error) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", err.ErrCode, err.ErrMsg)
}

// access_token 中控服务器接口.
type AccessTokenServer interface {
	Token() (token string, err error)                           // 请求中控服务器返回缓存的 access_token
	RefreshToken(currentToken string) (token string, err error) // 请求中控服务器刷新 access_token
	IID01332E16DF5011E5A9D5A4DB30FED8E1()                       // 接口标识, 没有实际意义
}

var _ AccessTokenServer = (*DefaultAccessTokenServer)(nil)

// DefaultAccessTokenServer 实现了 AccessTokenServer 接口.
//  NOTE:
//  1. 用于单进程环境.
//  2. 因为 DefaultAccessTokenServer 同时也是一个简单的中控服务器, 而不是仅仅实现 AccessTokenServer 接口,
//     所以整个系统只能存在一个 DefaultAccessTokenServer 实例!
type DefaultAccessTokenServer struct {
	appId     string
	appSecret string
	// httpClient *http.Client

	refreshTokenRequestChan  chan string             // chan currentToken
	refreshTokenResponseChan chan refreshTokenResult // chan {token, err}

	tokenCache unsafe.Pointer // *accessToken
}

// NewDefaultAccessTokenServer 创建一个新的 DefaultAccessTokenServer, 如果 httpClient == nil 则默认使用 util.DefaultHttpClient.
func NewDefaultAccessTokenServer(appId, appSecret string) (srv *DefaultAccessTokenServer) {

	srv = &DefaultAccessTokenServer{
		appId:     url.QueryEscape(appId),
		appSecret: url.QueryEscape(appSecret),
		// httpClient:               DefaultHttpClient,
		refreshTokenRequestChan:  make(chan string),
		refreshTokenResponseChan: make(chan refreshTokenResult),
	}

	go srv.tokenUpdateDaemon(time.Hour * 24 * time.Duration(100+rand.Int63n(200)))
	return
}

func (srv *DefaultAccessTokenServer) IID01332E16DF5011E5A9D5A4DB30FED8E1() {}

func (srv *DefaultAccessTokenServer) Token() (token string, err error) {
	if p := (*accessToken)(atomic.LoadPointer(&srv.tokenCache)); p != nil {
		return p.Token, nil
	}
	return srv.RefreshToken("")
}

type refreshTokenResult struct {
	token string
	err   error
}

func (srv *DefaultAccessTokenServer) RefreshToken(currentToken string) (token string, err error) {
	srv.refreshTokenRequestChan <- currentToken
	rslt := <-srv.refreshTokenResponseChan
	return rslt.token, rslt.err
}

func (srv *DefaultAccessTokenServer) tokenUpdateDaemon(initTickDuration time.Duration) {
	tickDuration := initTickDuration

NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case currentToken := <-srv.refreshTokenRequestChan:
			accessToken, cached, err := srv.updateToken(currentToken)
			if err != nil {
				srv.refreshTokenResponseChan <- refreshTokenResult{err: err}
				break
			}
			srv.refreshTokenResponseChan <- refreshTokenResult{token: accessToken.Token}
			if !cached {
				tickDuration = time.Duration(accessToken.ExpiresIn) * time.Second
				ticker.Stop()
				goto NEW_TICK_DURATION
			}

		case <-ticker.C:
			accessToken, _, err := srv.updateToken("")
			if err != nil {
				break
			}
			newTickDuration := time.Duration(accessToken.ExpiresIn) * time.Second
			if abs(tickDuration-newTickDuration) > time.Second*5 {
				tickDuration = newTickDuration
				ticker.Stop()
				goto NEW_TICK_DURATION
			}
		}
	}
}

func abs(x time.Duration) time.Duration {
	if x >= 0 {
		return x
	}
	return -x
}

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// updateToken 从微信服务器获取新的 access_token 并存入缓存, 同时返回该 access_token.
func (srv *DefaultAccessTokenServer) updateToken(currentToken string) (token *accessToken, cached bool, err error) {
	if currentToken != "" {
		if p := (*accessToken)(atomic.LoadPointer(&srv.tokenCache)); p != nil && currentToken != p.Token {
			return p, true, nil // 无需更改 p.ExpiresIn 参数值, cached == true 时用不到
		}
	}

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + srv.appId +
		"&secret=" + srv.appSecret

	httpResp, err := httpClient.Get(url)
	if err != nil {
		atomic.StorePointer(&srv.tokenCache, nil)
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		atomic.StorePointer(&srv.tokenCache, nil)
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result struct {
		Error
		accessToken
	}
	if err = DecodeJSONHttpResponse(httpResp.Body, &result); err != nil {
		atomic.StorePointer(&srv.tokenCache, nil)
		return
	}
	if result.ErrCode != ErrCodeOK {
		atomic.StorePointer(&srv.tokenCache, nil)
		err = &result.Error
		return
	}

	// 由于网络的延时, access_token 过期时间留有一个缓冲区
	switch {
	case result.ExpiresIn > 31556952: // 60*60*24*365.2425
		atomic.StorePointer(&srv.tokenCache, nil)
		err = errors.New("expires_in too large: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 10
	default:
		atomic.StorePointer(&srv.tokenCache, nil)
		err = errors.New("expires_in too small: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}

	tokenCopy := result.accessToken
	atomic.StorePointer(&srv.tokenCache, unsafe.Pointer(&tokenCopy))
	token = &tokenCopy
	return
}
