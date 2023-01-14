package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis/v9"
	"github.com/wangning057/scheduler/service/execute"
	"google.golang.org/grpc"
)

var taskReadyQueue = make(chan string, 1024)
var onExeCount = 0
var maxExeCount = 20

// 初始化一个redis客户端工具，以备使用
var rdb *redis.Client

func init() {
	fmt.Println("init in main.go")
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type executorServiceServer struct {
	execute.UnimplementedExecuteServiceServer
}

/*
	Execute 函数功能：

1.到redis中取任务
2.执行任务
*/
func (e *executorServiceServer) Execute(ctx context.Context, in *execute.ExecutionTask) (*execute.ExecuteResult, error) {
	// TODO
	action_id := in.GetActionId()
	oldStatus, err := rdb.GetSet(ctx, action_id, "done").Result()
	if err == redis.Nil {
		log.Fatalf("任务%+v在redis中不存在", action_id)
	} else  if err != nil {
		log.Fatalf("修改任务%+v在redis中的状态失败", action_id)
	}
	
	res := &execute.ExecuteResult{}

	if oldStatus != "ready" {
		return res, nil
	} else {
		 
	}



}

func main() {
	//相对于scheduler的gRPC服务端
	server := grpc.NewServer()
	execute.RegisterExecuteServiceServer(server, &executorServiceServer{})
	listener, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatal("从scheduler客户端接收到任务的服务监听端口失败", err)
	}
	_ = server.Serve(listener)
}
