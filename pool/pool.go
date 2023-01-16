package pool

import "github.com/wangning057/executor/runner"

//用来管理所有的 runner

type RunnerPool struct {
	RunnerCount int
}

func (p *RunnerPool)RunnersInitAndRun() {
	for i := 0; i < p.RunnerCount; i++ {
		newCmdRunner := &runner.CmdRunner{RunnerId: i}
		go newCmdRunner.Run()
	}
}