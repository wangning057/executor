package returner

import (
	"errors"
	"log"

	"github.com/wangning057/executor/service/execute"
)

//Returner负责在每条task执行前开启channel，
//将task_id与执行结果的channel对应
//

// 每一个task都有自己的一个channel
var resultMap = make(map[string]chan *execute.ExecuteResult, 1024)

func InitChan(task_id string) chan *execute.ExecuteResult {
	ch := make(chan *execute.ExecuteResult, 1)
	resultMap[task_id] = ch
	return ch
}

func SetRes(task_id string, res *execute.ExecuteResult) error {
	if ch, ok := resultMap[task_id]; ok {
		ch <- res
		return nil
	} else {
		return errors.New("找不到task_id:" + task_id + "对应的channel")
	}
}

func GetRes(task_id string) *execute.ExecuteResult {
	resChan := resultMap[task_id]
	if resChan == nil {
		log.Printf("resChan == nil，无法取到任务id=%v的res\n", task_id)
	}

	for {
		select {
		case res := <-resChan:
			delete(resultMap, task_id)
			close(resChan)
			return res
		}
	}
}
