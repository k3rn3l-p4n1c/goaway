package core

import (
	"errors"
	"github.com/k3rn3l-p4n1c/goaway/agent"
	"github.com/k3rn3l-p4n1c/goaway/file"
	"fmt"
	"github.com/k3rn3l-p4n1c/goaway/scheduler"
)

type Response struct {
	Message string
	Error bool
}

type Request struct {
	Command string
}

type Handler struct{}

func (h *Handler) GetName(req Request, res *Response) (err error) {
	if req.Command == "" {
		err = errors.New("A name must be specified")
		return
	}

	res.Message = "Hello " + req.Command
	return
}

var UpCmd = "Handler.Up"

func (h *Handler) Up(req Request, res *Response) (err error) {
	res.Message = "started"

	return
}

var GetCmd = "Handler.Get"

func (h *Handler) Get(req Request, res *Response) (err error) {
	v, _ := agent.GetInstance().ProjectName.Get()
	res.Message = v

	return
}

var SetCmd = "Handler.Set"

func (h *Handler) Set(req Request, res *Response) (err error) {
	agent.GetInstance().ProjectName.Set(req.Command)

	return
}

var StackUpCmd = "Handler.StackUp"

func (h *Handler) StackUp(req Request, res *Response) (err error) {
	stackFile, err := file.Read(req.Command)
	for i := range stackFile.Services {
		agent.GetInstance().Cluster.AddOrUpdateDeployment(i)
	}
	stack := scheduler.GenerateRandomStack(agent.GetInstance().Cluster)
	agent.GetInstance().Stack = &stack

	if err != nil {
		res.Error = true
		res.Message = string(err.Error())
		return
	}

	res.Message = fmt.Sprintf("%v", stackFile)
	return
}

var StackOptimizeCmd = "Handler.StackOptimize"

func (h *Handler) StackOptimize(req Request, res *Response) (err error) {
	res.Message = scheduler.Run(agent.GetInstance().Cluster, agent.GetInstance().Stack)
	return
}