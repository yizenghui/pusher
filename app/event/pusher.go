package event

import (
	"encoding/json"
	"fmt"
	"log"
	"pusher/common"
)

// 执行检查
func Add(uid, eid int, openid, name, keyindex string, cd, ov int, data map[string]string) error {

	var timekey = fmt.Sprintf("last_sent_timestamp_%d_%d", uid, eid) // 时间key
	var msgkey = fmt.Sprintf("last_messages_%d_%d", uid, eid)        // 内容key

	// 使用本地缓存存储用户最后一次发送站内信的时间戳（多服务器请用Redis或数据库）
	//	lastSentTimestamp, _ := common.GetCache(timekey, ``)

	lastMessagesStr, _ := common.GetCache(msgkey, ``)
	// // 聚合待推送消息
	// lastMessagesStr, _ := Cache("slmsg").Get(msgkey, []interface{}{})
	// lastMessages = append(lastMessages, map[string]interface{}{"title": title, "info": info, "username": "用户名", "date": date})

	var arr = []map[string]string{}

	err := json.Unmarshal([]byte(lastMessagesStr), &arr)
	if err != nil {
		log.Println(err)
	}

	log.Println(arr)

	log.Println(timekey, msgkey)
	//Collection Name
	return nil
}
