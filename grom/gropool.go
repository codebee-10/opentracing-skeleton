package main

import (
        "time"
        "fmt"
        "runtime"
)


type Task struct {
        function func() error
}


func NewTask(f func() error) *Task {
        t := Task{ function: f }
        return &t
}


func (t *Task) Execute() {
        t.function() //调用任务所绑定的函数
}


type Pool struct {
        TaskChannel chan *Task //对外接收Task的入口
        worker_num int  //协程池最大worker数量,限定Goroutine的个数
        JobsChannel chan *Task //协程池内部的任务就绪队列
}


func NewPool(cap int) *Pool {
        p := Pool{
                TaskChannel: make(chan *Task),
                worker_num:   cap,
                JobsChannel:  make(chan *Task),
        }

        return &p
}


func (p *Pool) worker(work_ID int) {
        for task := range p.JobsChannel {  //worker不断的从JobsChannel内部任务队列中拿任务
                task.Execute()  //如果拿到任务,则执行task任务
                fmt.Println("worker ID ", work_ID)
        }
}


func (p *Pool) Run() {
        //首先根据协程池的worker数量限定,开启固定数量的Worker
        for i := 0; i < p.worker_num; i++ {
                go p.worker(i)  //  每一个Worker用一个Goroutine承载
        }

        //从TaskChannel协程池入口取外界传递过来的任务
        //并且将任务送进JobsChannel中
        for task := range p.TaskChannel {
                p.JobsChannel <- task
                // time.Sleep(time.Second)
        }

        //3, 执行完毕需要关闭JobsChannel
        defer close(p.JobsChannel)

        //4, 执行完毕需要关闭TaskChannel
        defer close(p.TaskChannel)
}

//主函数
func main() {
        runtime.GOMAXPROCS(runtime.NumCPU())
        //创建一个Task
        t := NewTask(func() error {
                fmt.Println(time.Now())
                return nil
        })

        p := NewPool(3) //创建一个协程池,最大开启3个协程worker
        go func() {
                for i:=0;i<100;i++{
                        p.TaskChannel <- t  //开一个协程 不断的向 Pool 输送打印一条时间的task任务
                        fmt.Println(i, "......")
                }
        }()

       
        p.Run()  //启动协程池p

}




func producer(n int) <-chan int {
    out := make(chan int)
    go func() {
        defer func() {
            close(out)
            out = nil
            fmt.Println("producer exit")
        }()

        for i := 0; i < n*10; i++ {
            fmt.Printf("send %d\n", i)
            out <- i
            // time.Sleep(1 * time.Millisecond)
        }
    }()
    return out
}

// consumer only read data from in channel and print it
func consumer(in <-chan int) <-chan struct{} {
    finish := make(chan struct{})

    go func() {
        defer func() {
            fmt.Println("worker exit")
            finish <- struct{}{}
            close(finish)
        }()

        // Using for-range to exit goroutine
        // range has the ability to detect the close/end of a channel
        for x := range in {
            fmt.Printf("Process %d\n", x)
            // time.Sleep(1000 * time.Millisecond)
        }
    }()

    return finish
}

