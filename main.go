package main

import (
	"log"
	"net/http"
	"time"

	"github.com/yizenghui/pusher/app/job"

	//_ "D:\Go\src\github.com\yizenghui\pusher/docs" // 导入生成的docs.go文件

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type (
	message struct {
		OpenIDS     []string `json:"openids"`
		RID         string   `json:"rid"`
		TemplateID  string   `json:"template_id"`
		Page        string   `json:"page"`
		Key1        string   `json:"key1"`
		Value1      string   `json:"value1"`
		Key2        string   `json:"key2"`
		Value2      string   `json:"value2"`
		Key3        string   `json:"key3"`
		Value3      string   `json:"value3"`
		Key4        string   `json:"key4"`
		Value4      string   `json:"value4"`
		Key5        string   `json:"key5"`
		Value5      string   `json:"value5"`
		AccessToken string   `json:"access_token"`
		NoticeURL   string   `json:"notice_url"`
	}
)

type (
	template struct {
		OpenIDS             []string `json:"openids"`
		RID                 string   `json:"rid"`
		TemplateID          string   `json:"template_id"`
		PageURL             string   `json:"page_url"`
		MiniprogramAppID    string   `json:"miniprogram_appid"`
		MiniprogramPagePath string   `json:"miniprogram_page_path"`
		ClientMSGID         string   `json:"client_msg_id"`
		Key1                string   `json:"key1"`
		Value1              string   `json:"value1"`
		Color1              string   `json:"color1"`
		Key2                string   `json:"key2"`
		Value2              string   `json:"value2"`
		Color2              string   `json:"color2"`
		Key3                string   `json:"key3"`
		Value3              string   `json:"value3"`
		Color3              string   `json:"color3"`
		Key4                string   `json:"key4"`
		Value4              string   `json:"value4"`
		Color4              string   `json:"color4"`
		Key5                string   `json:"key5"`
		Value5              string   `json:"value5"`
		Color5              string   `json:"color5"`
		AccessToken         string   `json:"access_token"`
		NoticeURL           string   `json:"notice_url"`
	}
)

/*
*
'content-type': 'application/x-www-form-urlencoded'   c.FormValue 可以取得post的值
*/
func main() {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //  上海 Asia/Shanghai
	log.Println(time.Now().In(cstSh).Format("2006-01-02 15:04:05"))
	// if err := db.Init(); err != nil {
	// 	panic(fmt.Sprintf("mysql init failed with %+v", err))
	// }
	e := echo.New()
	// e.Use(middleware.RequestID())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello, World!")
	})
	e.POST("/msg", SendMessage)
	e.POST("/stl", SendTemplate)

	// 添加Swagger文档路由
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":1313")) //1323

}

// @Summary 小程序订阅消息推送
// @Description 获取指定用户的信息
// @ID getUserByID
// @Param id path int true "用户ID"
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/{id} [get]
func SendMessage(c echo.Context) error {

	msg := &message{}
	if err := c.Bind(msg); err != nil {
		return err
	}

	for _, openid := range msg.OpenIDS {
		// fmt.Println("openid", openid)
		plan := job.Message{
			OpenID:      openid,
			TemplateID:  msg.TemplateID,
			ReqID:       msg.RID,
			Page:        msg.Page,
			Key1:        msg.Key1,
			Value1:      msg.Value1,
			Key2:        msg.Key2,
			Value2:      msg.Value2,
			Key3:        msg.Key3,
			Value3:      msg.Value3,
			Key4:        msg.Key4,
			Value4:      msg.Value4,
			Key5:        msg.Key5,
			Value5:      msg.Value5,
			AccessToken: msg.AccessToken,
		}

		job.JobQueue <- job.Job{
			Player: &plan,
		}
		// log.Println(len(job.JobQueue))
	}

	return c.JSON(http.StatusOK, msg.RID)
}

/**
* 列队响应公众号模板消息推送事件
 */
func SendTemplate(c echo.Context) error {

	msg := &template{}
	if err := c.Bind(msg); err != nil {
		return err
	}

	for _, openid := range msg.OpenIDS {
		// fmt.Println("openid", openid)
		plan := job.Template{
			OpenID:              openid,
			TemplateID:          msg.TemplateID,
			ReqID:               msg.RID,
			PageURL:             msg.PageURL,
			MiniprogramAppID:    msg.MiniprogramAppID,
			MiniprogramPagePath: msg.MiniprogramPagePath,
			Key1:                msg.Key1,
			Value1:              msg.Value1,
			Color1:              msg.Color1,
			Key2:                msg.Key2,
			Value2:              msg.Value2,
			Color2:              msg.Color2,
			Key3:                msg.Key3,
			Value3:              msg.Value3,
			Color3:              msg.Color3,
			Key4:                msg.Key4,
			Value4:              msg.Value4,
			Color4:              msg.Color4,
			Key5:                msg.Key5,
			Value5:              msg.Value5,
			Color5:              msg.Color5,
			AccessToken:         msg.AccessToken,
		}

		job.JobQueue <- job.Job{
			Player: &plan,
		}
		// log.Println(len(job.JobQueue))
	}

	return c.JSON(http.StatusOK, msg.RID)
}

// func doSendMessage(uid int64, eid int64, title string, info string, username string, date string, time int, threshold int) bool {
// 	timekey := fmt.Sprintf("last_sent_timestamp_%d_%d", uid, eid) // 时间key
// 	msgkey := fmt.Sprintf("last_messages_%d_%d", uid, eid)        // 内容key

// 	// 使用本地缓存存储用户最后一次发送站内信的时间戳（多服务器请用Redis或数据库）
// 	lastSentTimestamp := cache.Get(timekey)

// 	// 聚合待推送消息
// 	lastMessages := cache.Get(msgkey)
// 	lastMessages = append(lastMessages, map[string]string{"title": title, "info": info, "username": username, "date": date})

// 	// 更新待推送消息
// 	cache.Put(msgkey, lastMessages, time*60+60)

// 	// 即时推送，需满足下列条件：没有上次推送时间，上次推送时间超过聚集时间，聚集信息量达到阈值
// 	if lastSentTimestamp == nil || time.Now().Unix()-lastSentTimestamp.(int64) > int64(time*60) || len(lastMessages) >= threshold {

// 		// 更新用户最后一次发送站内信的时间戳
// 		cache.Put(timekey, time.Now().Unix(), time*60+60)
// 		lastSentTimestamp = cache.Get(timekey)

// 		// 上次发送消息已经超过预定时限，将待发送的消息放入队列任务中进行处理
// 		SendJoinMessages(uid, eid)
// 		return true
// 	} else {
// 		// 更新用户最后一次消息的时间戳和聚合的消息列表
// 		cache.Put(msgkey, lastMessages, time*60+60)
// 		if len(lastMessages) == 1 { // 当只有一个消息要推送时，把推送任务安排下去
// 			// 把通知下发下去
// 			go SendJoinMessages(uid, eid)
// 			return true
// 		}
// 	}

// 	return false
// }

// func SendJoinMessages(userid, eid int64) {
// 	// 拿出当前缓存要推送的数据包
// 	messages := cache.Store("slmsg").Get(fmt.Sprintf("last_messages_%d_%d", userid, eid), []map[string]string{})
// 	if len(messages) > 0 {
// 		// 清除数据里面的记录
// 		cache.Store("slmsg").Put(fmt.Sprintf("last_messages_%d_%d", userid, eid), []map[string]string{})

// 		var username string
// 		message := messages[len(messages)-1]
// 		var usernames []string
// 		for _, m := range messages { // 整合数据包里面的信息
// 			usernames = append(usernames, m["username"])
// 		}
// 		username = strings.Join(usernames, "、")
// 		if len(usernames) > 1 { //
// 			username = username + "等" + strconv.Itoa(len(usernames)) + "人"
// 		}

// 		user := getUserByID(userid)
// 		if user != nil {
// 			params := map[string]interface{}{
// 				"thing2": message["info"],
// 				"thing3": username,
// 				"time1":  time.Now().Format("2006-01-02 15:04:05"),
// 			}
// 			templateID := config.Get("notice.media_solitaire_joined.template_id").(string)
// 			jumppage := "pages/relayDetail/relayDetail?id=" + strconv.FormatInt(eid, 10)
// 			user.notifyInternalMessage(params, templateID, jumppage)
// 		}
// 	}
// }
