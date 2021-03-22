// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: datamover.proto

package datamover

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Datamover service

type DatamoverService interface {
	Runjob(ctx context.Context, in *RunJobRequest, opts ...client.CallOption) (*RunJobResponse, error)
	DoLifecycleAction(ctx context.Context, in *LifecycleActionRequest, opts ...client.CallOption) (*LifecycleActionResonse, error)
	AbortJob(ctx context.Context, in *AbortJobRequest, opts ...client.CallOption) (*AbortJobResponse, error)
	PauseJob(ctx context.Context, in *PauseJobRequest, opts ...client.CallOption) (*PauseJobResponse, error)
}

type datamoverService struct {
	c    client.Client
	name string
}

func NewDatamoverService(name string, c client.Client) DatamoverService {
	return &datamoverService{
		c:    c,
		name: name,
	}
}

func (c *datamoverService) Runjob(ctx context.Context, in *RunJobRequest, opts ...client.CallOption) (*RunJobResponse, error) {
	req := c.c.NewRequest(c.name, "Datamover.Runjob", in)
	out := new(RunJobResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datamoverService) DoLifecycleAction(ctx context.Context, in *LifecycleActionRequest, opts ...client.CallOption) (*LifecycleActionResonse, error) {
	req := c.c.NewRequest(c.name, "Datamover.DoLifecycleAction", in)
	out := new(LifecycleActionResonse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datamoverService) AbortJob(ctx context.Context, in *AbortJobRequest, opts ...client.CallOption) (*AbortJobResponse, error) {
	req := c.c.NewRequest(c.name, "Datamover.AbortJob", in)
	out := new(AbortJobResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datamoverService) PauseJob(ctx context.Context, in *PauseJobRequest, opts ...client.CallOption) (*PauseJobResponse, error) {
	req := c.c.NewRequest(c.name, "Datamover.PauseJob", in)
	out := new(PauseJobResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Datamover service

type DatamoverHandler interface {
	Runjob(context.Context, *RunJobRequest, *RunJobResponse) error
	DoLifecycleAction(context.Context, *LifecycleActionRequest, *LifecycleActionResonse) error
	AbortJob(context.Context, *AbortJobRequest, *AbortJobResponse) error
	PauseJob(context.Context, *PauseJobRequest, *PauseJobResponse) error
}

func RegisterDatamoverHandler(s server.Server, hdlr DatamoverHandler, opts ...server.HandlerOption) error {
	type datamover interface {
		Runjob(ctx context.Context, in *RunJobRequest, out *RunJobResponse) error
		DoLifecycleAction(ctx context.Context, in *LifecycleActionRequest, out *LifecycleActionResonse) error
		AbortJob(ctx context.Context, in *AbortJobRequest, out *AbortJobResponse) error
		PauseJob(ctx context.Context, in *PauseJobRequest, out *PauseJobResponse) error
	}
	type Datamover struct {
		datamover
	}
	h := &datamoverHandler{hdlr}
	return s.Handle(s.NewHandler(&Datamover{h}, opts...))
}

type datamoverHandler struct {
	DatamoverHandler
}

func (h *datamoverHandler) Runjob(ctx context.Context, in *RunJobRequest, out *RunJobResponse) error {
	return h.DatamoverHandler.Runjob(ctx, in, out)
}

func (h *datamoverHandler) DoLifecycleAction(ctx context.Context, in *LifecycleActionRequest, out *LifecycleActionResonse) error {
	return h.DatamoverHandler.DoLifecycleAction(ctx, in, out)
}

func (h *datamoverHandler) AbortJob(ctx context.Context, in *AbortJobRequest, out *AbortJobResponse) error {
	return h.DatamoverHandler.AbortJob(ctx, in, out)
}

func (h *datamoverHandler) PauseJob(ctx context.Context, in *PauseJobRequest, out *PauseJobResponse) error {
	return h.DatamoverHandler.PauseJob(ctx, in, out)
}
