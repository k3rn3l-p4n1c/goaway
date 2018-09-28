package core

import (
	"errors"
	"github.com/k3rn3l-p4n1c/goaway/agent"
)

type Response struct {
	Message string
}

type Request struct {
	Name string
}

type Handler struct{}

func (h *Handler) GetName(req Request, res *Response) (err error) {
	if req.Name == "" {
		err = errors.New("A name must be specified")
		return
	}

	res.Message = "Hello " + req.Name
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
	println(v)

	return
}

var SetCmd = "Handler.Set"

func (h *Handler) Set(req Request, res *Response) (err error) {
	agent.GetInstance().ProjectName.Set(req.Name)

	return
}
