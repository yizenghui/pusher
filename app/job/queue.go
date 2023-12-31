package job

import (
	"fmt"
	"time"
)

// Player 任务接口
type Player interface {
	Run()
}

// PlayerDemo struct 示例结构
type PlayerDemo struct {
	OpenID string
	Text   string
}

// Run Player.Run
func (p *PlayerDemo) Run() {
	time.Sleep(time.Microsecond * 10)
	time.Sleep(time.Second)
	fmt.Printf("给 %s 发送内容 %s\n", p.OpenID, p.Text)
	time.Sleep(time.Second * 10)
	time.Sleep(time.Minute)
}

// Job represents the job to be run
type Job struct {
	Player Player
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				// if err := job.Payload.UploadToS3(); err != nil {
				// 	log.Errorf("Error uploading to S3: %s", err.Error())
				// }
				// fmt.Println(job)

				job.Player.Run()
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	MaxWorkers int
	WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, MaxWorkers: maxWorkers}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job

			}(job)
		}
	}
}

func init() {

	maxWorkers := 10 // 多少个并发
	maxQueue := 100  // 最多一瞬间可以入100个任务进行排队（够够的了）
	//初始化一个调试者,并指定它可以操作的 工人个数
	dispatch := NewDispatcher(maxWorkers)
	JobQueue = make(chan Job, maxQueue) //指定任务的队列长度
	//并让它一直接运行
	dispatch.Run()
	// close(notice.JobQueue)
}
