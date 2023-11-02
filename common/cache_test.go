// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func Test_Get(t *testing.T) {
	client := GetCacheClient()
	// if false {

	// 	t.Fatal(client)
	// }
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379", // Redis 服务器地址和端口号
	// 	Password: "",               // Redis 密码
	// 	DB:       0,                // Redis 数据库索引
	// })
	err := client.Set("key1", "value1", time.Minute).Err()
	if err != nil {
		log.Println("Error storing value in cache:", err)
	}

	// 从缓存中获取键对应的值
	val, err := client.Get("key1").Result()
	if err == redis.Nil { // 检查不存在的情况
		log.Println("Key does not exist in cache")
	} else if err != nil {
		log.Println("Error retrieving value from cache:", err)
	} else {
		log.Println(val)
	}

	// 删除缓存中的某个键值对
	err = client.Del("key1").Err()
	if err != nil {
		log.Println("Error deleting key from cache:", err)
	}
	t.Fatal(val)

}

func Test_Get22(t *testing.T) {
	client := GetCacheClient()
	// if false {

	// 	t.Fatal(client)
	// }
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379", // Redis 服务器地址和端口号
	// 	Password: "",               // Redis 密码
	// 	DB:       0,                // Redis 数据库索引
	// })

	array := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	jsonData, err := json.Marshal(array)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Set("key1", string(jsonData), time.Minute).Err()
	if err != nil {
		log.Println("Error storing value in cache:", err)
	}

	// 从缓存中获取键对应的值
	val, err := client.Get("key1").Result()
	if err == redis.Nil { // 检查不存在的情况
		log.Println("Key does not exist in cache")
	} else if err != nil {
		log.Println("Error retrieving value from cache:", err)
	} else {
		log.Println(val)
	}

	// 删除缓存中的某个键值对
	err = client.Del("key1").Err()
	if err != nil {
		log.Println("Error deleting key from cache:", err)
	}

	var cachedArray [][]int
	err = json.Unmarshal([]byte(val), &cachedArray)
	if err != nil {
		log.Fatal(err)
	}

	t.Fatal(cachedArray[0])
	t.Fatal(val)

}

func Test_Get33x(t *testing.T) {
	client := GetCacheClient()

	// array := [][]int{
	// 	{1, 2, 3},
	// 	{4, 5, 6},
	// 	{7, 8, 9},
	// }
	array := []int{1, 2, 3}

	err := client.Set("key3", array, time.Minute).Err()
	if err != nil {
		log.Println("Error storing value in cache:", err)

		t.Fatal(err)
	} else {

		// 从缓存中获取键对应的值
		val, err := client.Get("key3").Result()
		if err == redis.Nil { // 检查不存在的情况
			log.Println("Key does not exist in cache")
		} else if err != nil {
			log.Println("Error retrieving value from cache:", err)
		} else {
			log.Println(val)
		}

		var cachedArray [][]int
		err = json.Unmarshal([]byte(val), &cachedArray)
		if err != nil {
			log.Fatal(err)
		}

		t.Fatal(cachedArray[0])
		t.Fatal(val)
	}

}
