package event

import (
	"testing"
)

func Test_Add(t *testing.T) {
	//map[string]interface{}
	var data = map[string]string{"keyword1": "张三", "pagepath": "/pages/index/index"}
	err := Add(1, 1, `openid`, `username`, `keyword`, 1, 100, data)

	t.Fatal(err)
}
