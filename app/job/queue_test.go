// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

// func Test_T(t *testing.T) {

// 	//初始化对象池

// 	maxWorkers := 10
// 	maxQueue := 20
// 	//初始化一个调试者,并指定它可以操作的 工人个数
// 	dispatch := NewDispatcher(maxWorkers)
// 	JobQueue = make(chan Job, maxQueue) //指定任务的队列长度
// 	//并让它一直接运行
// 	dispatch.Run()

// 	for i := 0; i < 30; i++ {
// 		p := NoticeDemo{
// 			fmt.Sprintf("[%s]", strconv.Itoa(i)),
// 			fmt.Sprintf("玩家-[%s]", strconv.Itoa(i)),
// 		}
// 		JobQueue <- Job{
// 			Notice: &p,
// 		}
// 		log.Println(i, len(JobQueue))
// 		// time.Sleep(time.Millisecond)
// 	}
// 	time.Sleep(time.Second * 5)
// 	close(JobQueue)
// }

func Test_T3(t *testing.T) {

	for i := 0; i < 10; i++ {
		p := PlayerDemo{
			fmt.Sprintf("x[%s]", strconv.Itoa(i)),
			fmt.Sprintf("x玩家-[%s]", strconv.Itoa(i)),
		}
		JobQueue <- Job{
			Player: &p,
		}
		log.Println(i, len(JobQueue))
		// time.Sleep(time.Millisecond)
	}
	time.Sleep(time.Second * 3)
	log.Println(`JobQueue len `, len(JobQueue))
	if len(JobQueue) > 0 {
		close(JobQueue)
	}
}

// 	// https://blog.csdn.net/qq_31406415/article/details/110521553
// func Test_Tt(t *testing.T) {

// 	ch := make(chan int, 1) // 创建有缓冲channel
// 	go func() {
// 		fmt.Println("time sleep 5 second...")
// 		time.Sleep(5 * time.Second)
// 		<-ch
// 	}()
// 	ch <- 1 // 协程不会阻塞，等待数据被读取
// 	fmt.Println("第二次发送数据到channel， 即将阻塞")
// 	ch <- 1 // 第二次发送数据到channel, 在数据没有被读取之前，因为缓冲区满了， 所以会阻塞主协程。
// 	fmt.Println("ch 数据被消费，主协程退出")
// }

// func Test_Tnt(t *testing.T) {

// 	ch := make(chan int) // 创建无缓冲channel

// 	go func() {
// 		fmt.Println("time sleep 5 second...")
// 		time.Sleep(5 * time.Second)
// 		<-ch
// 	}()
// 	fmt.Println("即将阻塞...")
// 	ch <- 1 // 协程将会阻塞，等待数据被读取
// 	fmt.Println("ch 数据被消费，主协程退出")
// }
