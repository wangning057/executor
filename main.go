package main

import (
	"context"
	"log"
	"net"
	"runtime"

	"github.com/wangning057/executor/pool"
	q "github.com/wangning057/executor/queue"
	"github.com/wangning057/executor/returner"
	"github.com/wangning057/executor/service/execute"
	"google.golang.org/grpc"
)

type executorServiceServer struct {
	execute.UnimplementedExecuteServiceServer
}

/*
	Execute 函数功能：

1.到redis中取任务
2.执行任务
3.返回结果
*/
func (e *executorServiceServer) Execute(ctx context.Context, task *execute.ExecutionTask) (*execute.ExecuteResult, error) {
	// TODO
	task_id := task.GetTaskId()
	q.ReadyQueue <- task
	res := returner.GetRes(task_id)
	return res, nil
}

func main() {

	//获取机器cpu核数
	cpuCount := runtime.NumCPU()

	//初始化一个runnerPool，并启动
	runnerPool := &pool.RunnerPool{
		RunnerCount: cpuCount + 2,
	}
	runnerPool.RunnersInitAndRun()

	//相对于scheduler的gRPC服务端
	server := grpc.NewServer()
	execute.RegisterExecuteServiceServer(server, &executorServiceServer{})
	listener, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatal("从scheduler客户端接收到任务的服务监听端口失败", err)
	}
	_ = server.Serve(listener)
}
