// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/example/example.proto

/*
Package go_micro_srv_data is a generated protocol buffer package.

It is generated from these files:
	proto/example/example.proto

It has these top-level messages:
	Request
	Response
	LoginLampRequest
	LoginLampResponse
	Data
	ConfigLampRequest
	ConfigLampResponse
	New
	Lamp
	PullNewLampRequest
	PullNewLampResponse
	GoldLampRequest
	GoldLampResponse
	RankingLampRequest
	RankingLampResponse
	CheckpointLampRequest
	CheckpointLampResponse
	SetLampRequest
	SetLampResponse
	GetLampRequest
	GetLampResponse
	BuyLampRequest
	BuyLampResponse
*/
package go_micro_srv_data

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Example service

type ExampleService interface {
	Data(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// 登陆服务
	LoginLamp(ctx context.Context, in *LoginLampRequest, opts ...client.CallOption) (*LoginLampResponse, error)
	// 拉取邀请新人列表服务
	PullNewLamp(ctx context.Context, in *PullNewLampRequest, opts ...client.CallOption) (*PullNewLampResponse, error)
	// 初始化服务
	ConfigLamp(ctx context.Context, in *ConfigLampRequest, opts ...client.CallOption) (*ConfigLampResponse, error)
	// 金币变化服务
	GoldLamp(ctx context.Context, in *GoldLampRequest, opts ...client.CallOption) (*GoldLampResponse, error)
	// 服务器排名服务
	RankingLamp(ctx context.Context, in *RankingLampRequest, opts ...client.CallOption) (*RankingLampResponse, error)
	// 关卡服务
	CheckpointLamp(ctx context.Context, in *CheckpointLampRequest, opts ...client.CallOption) (*CheckpointLampResponse, error)
	// 设置灯服务
	SetLamp(ctx context.Context, in *SetLampRequest, opts ...client.CallOption) (*SetLampResponse, error)
	// 获取灯列表服务
	GetLamp(ctx context.Context, in *GetLampRequest, opts ...client.CallOption) (*GetLampResponse, error)
	// 购买灯服务
	BuyLamp(ctx context.Context, in *BuyLampRequest, opts ...client.CallOption) (*BuyLampResponse, error)
}

type exampleService struct {
	c    client.Client
	name string
}

func NewExampleService(name string, c client.Client) ExampleService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.data"
	}
	return &exampleService{
		c:    c,
		name: name,
	}
}

func (c *exampleService) Data(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Example.Data", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) LoginLamp(ctx context.Context, in *LoginLampRequest, opts ...client.CallOption) (*LoginLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.LoginLamp", in)
	out := new(LoginLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) PullNewLamp(ctx context.Context, in *PullNewLampRequest, opts ...client.CallOption) (*PullNewLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.PullNewLamp", in)
	out := new(PullNewLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) ConfigLamp(ctx context.Context, in *ConfigLampRequest, opts ...client.CallOption) (*ConfigLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.ConfigLamp", in)
	out := new(ConfigLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) GoldLamp(ctx context.Context, in *GoldLampRequest, opts ...client.CallOption) (*GoldLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.GoldLamp", in)
	out := new(GoldLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) RankingLamp(ctx context.Context, in *RankingLampRequest, opts ...client.CallOption) (*RankingLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.RankingLamp", in)
	out := new(RankingLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) CheckpointLamp(ctx context.Context, in *CheckpointLampRequest, opts ...client.CallOption) (*CheckpointLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.CheckpointLamp", in)
	out := new(CheckpointLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) SetLamp(ctx context.Context, in *SetLampRequest, opts ...client.CallOption) (*SetLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.SetLamp", in)
	out := new(SetLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) GetLamp(ctx context.Context, in *GetLampRequest, opts ...client.CallOption) (*GetLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.GetLamp", in)
	out := new(GetLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) BuyLamp(ctx context.Context, in *BuyLampRequest, opts ...client.CallOption) (*BuyLampResponse, error) {
	req := c.c.NewRequest(c.name, "Example.BuyLamp", in)
	out := new(BuyLampResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Example service

type ExampleHandler interface {
	Data(context.Context, *Request, *Response) error
	// 登陆服务
	LoginLamp(context.Context, *LoginLampRequest, *LoginLampResponse) error
	// 拉取邀请新人列表服务
	PullNewLamp(context.Context, *PullNewLampRequest, *PullNewLampResponse) error
	// 初始化服务
	ConfigLamp(context.Context, *ConfigLampRequest, *ConfigLampResponse) error
	// 金币变化服务
	GoldLamp(context.Context, *GoldLampRequest, *GoldLampResponse) error
	// 服务器排名服务
	RankingLamp(context.Context, *RankingLampRequest, *RankingLampResponse) error
	// 关卡服务
	CheckpointLamp(context.Context, *CheckpointLampRequest, *CheckpointLampResponse) error
	// 设置灯服务
	SetLamp(context.Context, *SetLampRequest, *SetLampResponse) error
	// 获取灯列表服务
	GetLamp(context.Context, *GetLampRequest, *GetLampResponse) error
	// 购买灯服务
	BuyLamp(context.Context, *BuyLampRequest, *BuyLampResponse) error
}

func RegisterExampleHandler(s server.Server, hdlr ExampleHandler, opts ...server.HandlerOption) error {
	type example interface {
		Data(ctx context.Context, in *Request, out *Response) error
		LoginLamp(ctx context.Context, in *LoginLampRequest, out *LoginLampResponse) error
		PullNewLamp(ctx context.Context, in *PullNewLampRequest, out *PullNewLampResponse) error
		ConfigLamp(ctx context.Context, in *ConfigLampRequest, out *ConfigLampResponse) error
		GoldLamp(ctx context.Context, in *GoldLampRequest, out *GoldLampResponse) error
		RankingLamp(ctx context.Context, in *RankingLampRequest, out *RankingLampResponse) error
		CheckpointLamp(ctx context.Context, in *CheckpointLampRequest, out *CheckpointLampResponse) error
		SetLamp(ctx context.Context, in *SetLampRequest, out *SetLampResponse) error
		GetLamp(ctx context.Context, in *GetLampRequest, out *GetLampResponse) error
		BuyLamp(ctx context.Context, in *BuyLampRequest, out *BuyLampResponse) error
	}
	type Example struct {
		example
	}
	h := &exampleHandler{hdlr}
	return s.Handle(s.NewHandler(&Example{h}, opts...))
}

type exampleHandler struct {
	ExampleHandler
}

func (h *exampleHandler) Data(ctx context.Context, in *Request, out *Response) error {
	return h.ExampleHandler.Data(ctx, in, out)
}

func (h *exampleHandler) LoginLamp(ctx context.Context, in *LoginLampRequest, out *LoginLampResponse) error {
	return h.ExampleHandler.LoginLamp(ctx, in, out)
}

func (h *exampleHandler) PullNewLamp(ctx context.Context, in *PullNewLampRequest, out *PullNewLampResponse) error {
	return h.ExampleHandler.PullNewLamp(ctx, in, out)
}

func (h *exampleHandler) ConfigLamp(ctx context.Context, in *ConfigLampRequest, out *ConfigLampResponse) error {
	return h.ExampleHandler.ConfigLamp(ctx, in, out)
}

func (h *exampleHandler) GoldLamp(ctx context.Context, in *GoldLampRequest, out *GoldLampResponse) error {
	return h.ExampleHandler.GoldLamp(ctx, in, out)
}

func (h *exampleHandler) RankingLamp(ctx context.Context, in *RankingLampRequest, out *RankingLampResponse) error {
	return h.ExampleHandler.RankingLamp(ctx, in, out)
}

func (h *exampleHandler) CheckpointLamp(ctx context.Context, in *CheckpointLampRequest, out *CheckpointLampResponse) error {
	return h.ExampleHandler.CheckpointLamp(ctx, in, out)
}

func (h *exampleHandler) SetLamp(ctx context.Context, in *SetLampRequest, out *SetLampResponse) error {
	return h.ExampleHandler.SetLamp(ctx, in, out)
}

func (h *exampleHandler) GetLamp(ctx context.Context, in *GetLampRequest, out *GetLampResponse) error {
	return h.ExampleHandler.GetLamp(ctx, in, out)
}

func (h *exampleHandler) BuyLamp(ctx context.Context, in *BuyLampRequest, out *BuyLampResponse) error {
	return h.ExampleHandler.BuyLamp(ctx, in, out)
}
