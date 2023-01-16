package queue

import "github.com/wangning057/executor/service/execute"

var ReadyQueue = make(chan *execute.ExecutionTask, 1024)