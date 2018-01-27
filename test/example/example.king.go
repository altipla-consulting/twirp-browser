package example

import (
	"github.com/altipla-consulting/king/runtime"
	common "github.com/altipla-consulting/king/test/common"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Code generated by protoc-gen-king_go 1.0.0, DO NOT EDIT.
// Source: test/example/example.proto

type ContactMessagesServiceServer interface {
	Foo(ctx context.Context, in *FooRequest) (out *common.Empty, err error)
	Bar(ctx context.Context, in *BarRequest) (out *common.Empty, err error)
}

func RegisterContactMessagesService(server ContactMessagesServiceServer) {
	serviceDef := &runtime.Service{
		Name: "king.example.ContactMessagesService",
		Methods: []*runtime.Method{

			{
				Name:  "Foo",
				Input: func() proto.Message { return new(FooRequest) },
				Handler: func(ctx context.Context, in proto.Message) (proto.Message, error) {
					return server.Foo(ctx, in.(*FooRequest))
				},
			},

			{
				Name:  "Bar",
				Input: func() proto.Message { return new(BarRequest) },
				Handler: func(ctx context.Context, in proto.Message) (proto.Message, error) {
					return server.Bar(ctx, in.(*BarRequest))
				},
			},
		},
	}
	runtime.Services = append(runtime.Services, serviceDef)
}

type ContactMessagesServiceClient interface {
	Foo(ctx context.Context, in *FooRequest) (out *common.Empty, err error)
	Bar(ctx context.Context, in *BarRequest) (out *common.Empty, err error)
}

type clientImplContactMessagesService struct {
	server string
}

func NewContactMessagesServiceClient(server string) ContactMessagesServiceClient {
	return &clientImplContactMessagesService{server}
}

func (impl *clientImplContactMessagesService) Foo(ctx context.Context, in *FooRequest) (out *common.Empty, err error) {
	out = new(common.Empty)
	if err := runtime.ClientCall(ctx, impl.server, "king.example.ContactMessagesService", "Foo", in, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (impl *clientImplContactMessagesService) Bar(ctx context.Context, in *BarRequest) (out *common.Empty, err error) {
	out = new(common.Empty)
	if err := runtime.ClientCall(ctx, impl.server, "king.example.ContactMessagesService", "Bar", in, out); err != nil {
		return nil, err
	}

	return out, nil
}
