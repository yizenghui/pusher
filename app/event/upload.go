package event

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/yizenghui/pusher/app/mp"
	"github.com/yizenghui/pusher/db"
)

// 执行检查
func RunCheck(fileID int, hash string, fileSize int, token string) error {
	DB := db.Get() //.Debug()

	file := &db.File{}
	user := &db.User{}
	touser := &db.User{}

	// 获取文件参数
	DB.Table(`file_upload_record`).Find(&file, fileID)

	if file.ID == 0 {
		// log.Println(fileID)
		fmt.Println(fmt.Sprintf("not found id[%s]", strconv.Itoa(fileID)))
		return nil
	} else if file.CheckAt.IsZero() { //检查时间为空时进行送检逻辑
		// t.Fatal(time.Now().Format("2006-01-02 15:04:05"))
		//获取用户参数
		DB.Table(`users`).Select([]string{"openid", "posted", "related", "fromid", "isblock"}).Find(&user, file.UserID)
		if file.UserID != file.ToUserID { //给别人才有，给自己没有
			DB.Table(`users`).Select([]string{"openid", "posted", "related", "fromid", "isblock"}).Find(&touser, file.ToUserID)
		}
	} else { //已经检查了，没必要进来
		// log.Println(file)
		fmt.Println(fmt.Sprintf("resubmit id[%s]", strconv.Itoa(fileID)))
		return nil
	}

	if touser.Related < 6 { //送检（用户自传或接收人关联较少）
		// t.Fatal(touser)
		if file.MediaCate == 1 { //这是图片类型

			var m = mp.MediaSecCheck(file.URL, 2, user.Openid, token)
			if traceID, ok := m["trace_id"]; ok {
				ft := db.FileTrace{UserID: file.UserID, FileID: file.ID, TraceID: traceID.(string)}
				DB.Create(&ft) // 插入数据
				// fmt.Println(fmt.Sprintf("secci [%s]", strconv.Itoa(fileID)))
			} else {
				fmt.Println(fmt.Sprintf("seccfi [%s]", strconv.Itoa(fileID)))
				log.Println(m) // 这里是有问题的
				return nil
			}
		} else if file.MediaCate == 3 { //音频

			var m = mp.MediaSecCheck(file.URL, 1, user.Openid, token)
			if traceID, ok := m["trace_id"]; ok {
				ft := db.FileTrace{UserID: file.UserID, FileID: file.ID, TraceID: traceID.(string)}
				DB.Create(&ft) // 插入数据
				// fmt.Println(fmt.Sprintf("secca [%s]", strconv.Itoa(fileID)))
			} else {
				fmt.Println(fmt.Sprintf("seccfa [%s]", strconv.Itoa(fileID)))
				log.Println(m) // 这里是有问题的
				return nil
			}

		} else if file.MediaCate == 2 { //todo 视频

		} else { //todo 其他文件 注意 5 unkown (改作jpg图片)

		}
		if file.CompleteAt.IsZero() { //如果没有标记提交完成时间，则标记上
			DB.Table(`file_upload_record`).Where("id = ?", file.ID).Updates(
				map[string]interface{}{"hash": hash, "file_size": fileSize, "complete_at": time.Now().Format("2006-01-02 15:04:05")},
			)
		}
		// log.Println(file)
		return nil
	} else { // 免检
		if file.CompleteAt.IsZero() { //如果没有标记提交完成时间，则标记上
			DB.Table(`file_upload_record`).Where("id = ?", file.ID).Updates(
				map[string]interface{}{"hash": hash, "status": 5, "file_size": fileSize, "complete_at": time.Now().Format("2006-01-02 15:04:05")},
			)
		}
		// fmt.Println(fmt.Sprintf("complete id[%s]", strconv.Itoa(fileID)))
	}
	// log.Println(file)
	return nil
}
