package execute

import context "context"

var ExecuteService = &executeService{}

// queueService 用于实现接口 ExecuteServiceServer
type executeService struct {
}

/* Execute 函数功能：
1.从scheduler接收任务
2.查询redis中该任务的状态
3.改变redis中该任务的状态
4.将接收到的任务放入自己的待执行队列中
5.将完成队列中的执行结果发送给scheduler
*/
func (e *executeService) Execute(context.Context, *Command) (*ExecuteResult, error) {
	// TODO
}
func (e *executeService) mustEmbedUnimplementedEnQueueServiceServer() {}