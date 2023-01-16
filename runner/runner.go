package runner

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/wangning057/executor/cmd_util"
	q "github.com/wangning057/executor/queue"
	"github.com/wangning057/executor/returner"
	"github.com/wangning057/executor/service/execute"
)

/*设计在该 executor 初始化时，启动 maxExeCount 个 cmdRunner，每个 cmdRunner 是一个协程
将 readyQueue 初始化为 channel 类型
让每个 cmdRunner 自己到 readyQueue 中去不断地 “取一条任务，执行一条任务，取一条任务，执行一条任务。。”
每个 cmdRunner 也相当于是在抢任务

这样就可以完美控制并发执行的任务数！！

先搞一个 readyQueue ，以后再改进。让IO型任务有更高的优先级，似乎意义不大，因为并行的任务数是固定的。
让更慢的IO型任务先执行，反而会导致单位时间内任务完成数量减少。
正确的做法是优先将IO型任务分配给本地的executor执行
*/

const (
	WORK_DIR string = "/home/ubuntu/ninja" //命令的执行目录
)

// 初始化一个redis客户端工具，以备使用
var rdb *redis.Client

func init() {
	fmt.Println("init in exeFunc.go")
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type CmdRunner struct {
	RunnerId int
}

func (c *CmdRunner) Run() {
	c.initRunner()
	for {
		newTask := c.getTask()
		c.exeTask(newTask)
	}
}

func (c *CmdRunner) initRunner() {
	//进行runner的环境设置等初始化操作
	fmt.Println("runner", c.RunnerId, "初始化")
	c.Execute_Test()
}

func (c *CmdRunner) getTask() *execute.ExecutionTask {
	newTask := <-q.ReadyQueue
	return newTask
}

func (c *CmdRunner) exeTask(newTask *execute.ExecutionTask) error {
	taskId := newTask.GetTaskId()
	oldStatus, err := rdb.GetSet(context.Background(), taskId, "done").Result()
	if err == redis.Nil {
		log.Fatalf("任务%+v在redis中不存在", taskId)
	} else if err != nil {
		log.Fatalf("修改任务%+v在redis中的状态失败", taskId)
	}

	if oldStatus != "ready" {
		//略过该条任务
		//直接返回nil代表指令未执行 //这里不对了，应该是不需要通过channel返回
		return nil
	} else {
		//执行该条任务，初始化map中对应的channel
		returner.InitChan(taskId)

		ctx := context.Background()
		res := RunTask(ctx, newTask)

		//任务执行结束，将结果返回到 map 中的对应的channel中
		err2 := returner.SetRes(taskId, res)

		return err2
	}
}

// 执行传入的任务 t
func RunTask(ctx context.Context, t *execute.ExecutionTask) *execute.ExecuteResult {
	cmd := &cmd_util.Command{
		Content:     t.GetCommand(),
		Env:         make([]string, 0),
		Use_console: false,
	}
	var stdout, stderr bytes.Buffer
	stdio := &cmd_util.Stdio{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	cmd_util.Run(ctx, cmd, WORK_DIR, stdio)
	ctx.Done()
	res := &execute.ExecuteResult{}
	res.Signal = "ok"
	return res
}

// 加上(c *CmdRunner)
func (c *CmdRunner) Execute_Test() {

	ctx := context.Background()

	cmd := &cmd_util.Command{
		Content:     "pwd",
		Env:         make([]string, 0),
		Use_console: false,
	}

	var stdout, stderr bytes.Buffer
	stdio := &cmd_util.Stdio{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	fmt.Print(c.RunnerId, "测试：")
	cmd_util.Run(ctx, cmd, WORK_DIR, stdio)

	fmt.Println(c.RunnerId, "测试：", stdio.Stdout)

}
