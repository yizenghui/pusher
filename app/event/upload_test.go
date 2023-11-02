package event

import (
	"encoding/json"
	"fmt"
	"pusher/app/mp"
	"pusher/db"
	"testing"
	"time"
)

func Test_GetUser(t *testing.T) {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}
	DB := db.Get()

	// DB.AutoMigrate(&db.Project{})

	user := &db.User{}
	// DB.Debug().Find(&user, datatypes.JSONQuery("attributes").HasKey("role"))
	DB.Debug().Table(`users`).Select([]string{"openid", "posted", "related", "fromid", "isblock"}).Find(&user, 1) //
	// DB.Debug().Find(&user, datatypes.JSONQuery("attributes").HasKey("orgs", "orga"))

	// DB.Debug().First(&user, datatypes.JSONQuery("attributes").Equals("jinzhu", "name"))
	//	t.Fatal(user)

	file := &db.File{}
	DB.Debug().Table(`file_upload_record`).Find(&file, 6410299) //.Select([]string{"openid", "posted", "related", "fromid", "isblock"})
	// Where("id IN ?", []int{10, 11}).
	DB.Debug().Table(`file_upload_record`).Where("id = ?", file.ID).Updates(
		map[string]interface{}{"hash": "hellohash", "complete_at": "2022-11-04 16:10:20"},
	)
	// file.Hash = "jinzhu 2"
	// file.CompleteAt = time.Now()
	// DB.Debug().Save(&file)

	t.Fatal(file)

}

func Test_RunEventJob(t *testing.T) {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	var fileID = 6410299
	var token = `62_LMf3mS0lunDgAJALCN`
	var hash = `hellohash`
	var fileSize = 612800

	t.Fatal(RunCheck(fileID, hash, fileSize, token))
}

func Test_RunJob(t *testing.T) {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	var fileID = 6410299
	var token = `62_LMf3mS0lunprUUfKWpuOlV9wEW2ythEyLjhmyXgohca66kvqfGbD6av_JRn2FL5eRvfz5KzuQk2JLAjJwJTpP4F2kjGxcTyVLbG0ryCj1D12ZwUeioaMsqZE5pb_vU9cthX-k9fcwxEI96GsBSDgAJALCN`
	var hash = `hellohash`
	var fileSize = 612800

	t.Fatal(RunCheck(fileID, hash, fileSize, token))
	DB := db.Get()

	file := &db.File{}
	user := &db.User{}
	touser := &db.User{}

	// 获取文件参数
	DB.Debug().Table(`file_upload_record`).Find(&file, fileID)

	if file.CheckAt.IsZero() { //检查时间为空时进行送检逻辑
		// t.Fatal(time.Now().Format("2006-01-02 15:04:05"))
		//获取用户参数
		DB.Debug().Table(`users`).Select([]string{"openid", "posted", "related", "fromid", "isblock"}).Find(&user, file.UserID)
		if file.UserID != file.ToUserID { //给别人才有，给自己没有
			DB.Debug().Table(`users`).Select([]string{"openid", "posted", "related", "fromid", "isblock"}).Find(&touser, file.ToUserID)
		}
	} else { //已经检查了，没必要进来
		t.Fatal(file)
	}

	if touser.Related < 6 { //送检（用户自传或接收人关联较少）
		// t.Fatal(touser)
		if file.MediaCate == 1 { //这是图片类型

			var m = mp.MediaSecCheck(file.URL, 2, user.Openid, token)
			if traceID, ok := m["trace_id"]; ok {
				ft := db.FileTrace{UserID: file.UserID, FileID: file.ID, TraceID: traceID.(string)}
				DB.Debug().Create(&ft) // 插入数据
			}
			t.Fatal(m)

		} else if file.MediaCate == 3 { //音频

			var m = mp.MediaSecCheck(file.URL, 1, user.Openid, token)
			if traceID, ok := m["trace_id"]; ok {
				ft := db.FileTrace{UserID: file.UserID, FileID: file.ID, TraceID: traceID.(string)}
				DB.Debug().Create(&ft) // 插入数据
			}
			t.Fatal(m)

		} else if file.MediaCate == 2 { //todo 视频

		} else { //todo 其他文件 注意 5 unkown (改作jpg图片)

		}

		t.Fatal(file)

	} else { // 免检
		if file.CompleteAt.IsZero() { //如果没有标记提交完成时间，则标记上
			DB.Debug().Table(`file_upload_record`).Where("id = ?", file.ID).Updates(
				map[string]interface{}{"hash": hash, "status": 5, "file_size": fileSize, "complete_at": time.Now().Format("2006-01-02 15:04:05")},
			)
		}
	}
	t.Fatal(file)

}

func Test_PostCheck(t *testing.T) {

	token := `62_m63wHbSYpNborGenBIVXSnJwxeKk1kAJJcyY8L68fOF_TTePorXkFWRKUxFEUcEMmwEQhdWTZB04aLP1KdwV__ffVKUkiX1164Tk3-I6bzHkSFhs2pkLZNpObayQR239PtLLSNpp3RtjDQR6OUNcAEAZFY`
	urlStr := `https://api.weixin.qq.com/wxa/media_check_async?access_token=` + token

	// 	data := `{"content":"特3456书 yuuo 莞6543李 zxcz 蒜7782法 fgnv 级
	// 	完2347全 dfji 试3726测 asad 感3847知 qwez 到
	// "}` //习近平

	// 详情看官方文档
	type MediaCheckData struct {
		MediaURL  string `json:"media_url"`
		MediaType int32  `json:"media_type"` //  1:音频;2:图片
		Version   int32  `json:"version"`    //2
		Openid    string `json:"openid"`
		Scene     int32  `json:"scene"` // 2
	}

	// data := &Data{
	// 	Content: `特3456书 yuuo 莞6543李 zxcz 蒜7782法 fgnv 级
	// 完2347全 dfji 试3726测 asad 感3847知 qwez 到`,
	// }

	data := &MediaCheckData{
		MediaURL:  `https://tpc.googlesyndication.com/simgad/1551983530139128296?sqp=4sqPyQQrQikqJwhfEAEdAAC0QiABKAEwCTgDQPCTCUgAUAFYAWBfcAJ4AcUBLbKdPg&rs=AOga4qnQjYgguexVMcEzWblX8DA7q3uCBA`,
		MediaType: 2,
		Version:   2,
		Openid:    `221asdasd`,
		Scene:     2,
	}

	datajson, _ := json.Marshal(data)

	// t.Fatal(datajson)
	// source = fmt.Sprintf(source, user, pwd, addr, dataBase)
	// resp, err := http.Post(urlStr, "application/json", strings.NewReader(data))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// m, ok := gjson.ParseBytes([]byte(body)).Value().(map[string]interface{})
	// log.Println(m["name"], ok)
	// // reader, err := charset.NewReader(resp.Body, strings.ToLower(resp.Header.Get("Content-Type")))
	// defer resp.Body.Close()

	var m = mp.CloudHttpPost(urlStr, datajson)

	t.Fatal(m)
	// token 62_IX-KQN0_qDWeKGtKS5-n16C04ood31yM5u-zJXIFpKpa8E3VG7p9aDFwoaS2mnne9z_nBGgIPGj9Dm5wUiH6-yHr3-RxPjhjXoT45d_zc21NDHclfOnMZUctLsIGSjy7MHY07w3OIyevfvMIPMVdAGAEBK
	// mp.HTTPPostJSON('https://api.weixin.qq.com/wxa/media_check_async',)

	// $result = cloudHttpPost(  [
	// 	'media_url'=> $file->url,
	// 	'media_type' => $file->media_cate == 1?2:1, // 1:音频;2:图片
	// 	'version' => 2,
	// 	'openid' => $file->user->openid,
	// 	'scene' => 2,
	// ]);

	// if ( $result['errcode'] == 0 && $result['errmsg'] == 'ok' ){
	// 	$file_trace = new FileTrace;
	// 	$file_trace->file_id = $file->id;
	// 	$file_trace->user_id = $file->user_id;
	// 	$file_trace->trace_id = $result['trace_id'];
	// 	$file_trace->save();
	// 	// \Cache::put('mediacheck.'.$result['trace_id'], $file->id, 60);
	// }else{
	// 	\Log::channel('toadmincheck')->info('#'.$file->id.' res: '.var_export($result,true) );
	// }
}

func Test_TPostData(t *testing.T) {

	token := `62_m63wHbSYpNborGenBIVXSnJwxeKk1kAJJcyY8L68fOF_TTePorXkFWRKUxFEUcEMmwEQhdWTZB04aLP1KdwV__ffVKUkiX1164Tk3-I6bzHkSFhs2pkLZNpObayQR239PtLLSNpp3RtjDQR6OUNcAEAZFY`
	urlStr := `https://api.weixin.qq.com/wxa/msg_sec_check?access_token=` + token

	// 	data := `{"content":"特3456书 yuuo 莞6543李 zxcz 蒜7782法 fgnv 级
	// 	完2347全 dfji 试3726测 asad 感3847知 qwez 到
	// "}` //习近平

	type Data struct {
		Content string `json:"content"`
	}

	// data := &Data{
	// 	Content: `特3456书 yuuo 莞6543李 zxcz 蒜7782法 fgnv 级
	// 完2347全 dfji 试3726测 asad 感3847知 qwez 到`,
	// }

	data := &Data{
		Content: `免费成人网站`,
	}

	datajson, _ := json.Marshal(data)

	// t.Fatal(datajson)
	// source = fmt.Sprintf(source, user, pwd, addr, dataBase)
	// resp, err := http.Post(urlStr, "application/json", strings.NewReader(data))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// m, ok := gjson.ParseBytes([]byte(body)).Value().(map[string]interface{})
	// log.Println(m["name"], ok)
	// // reader, err := charset.NewReader(resp.Body, strings.ToLower(resp.Header.Get("Content-Type")))
	// defer resp.Body.Close()

	var m = mp.CloudHttpPost(urlStr, datajson)

	t.Fatal(m)
	// token 62_IX-KQN0_qDWeKGtKS5-n16C04ood31yM5u-zJXIFpKpa8E3VG7p9aDFwoaS2mnne9z_nBGgIPGj9Dm5wUiH6-yHr3-RxPjhjXoT45d_zc21NDHclfOnMZUctLsIGSjy7MHY07w3OIyevfvMIPMVdAGAEBK
	// mp.HTTPPostJSON('https://api.weixin.qq.com/wxa/media_check_async',)

	// $result = cloudHttpPost(  [
	// 	'media_url'=> $file->url,
	// 	'media_type' => $file->media_cate == 1?2:1, // 1:音频;2:图片
	// 	'version' => 2,
	// 	'openid' => $file->user->openid,
	// 	'scene' => 2,
	// ]);

	// if ( $result['errcode'] == 0 && $result['errmsg'] == 'ok' ){
	// 	$file_trace = new FileTrace;
	// 	$file_trace->file_id = $file->id;
	// 	$file_trace->user_id = $file->user_id;
	// 	$file_trace->trace_id = $result['trace_id'];
	// 	$file_trace->save();
	// 	// \Cache::put('mediacheck.'.$result['trace_id'], $file->id, 60);
	// }else{
	// 	\Log::channel('toadmincheck')->info('#'.$file->id.' res: '.var_export($result,true) );
	// }
}
