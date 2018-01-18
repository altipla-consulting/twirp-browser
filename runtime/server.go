package runtime

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

var Services []*Service

type Service struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name    string
	Input   func() proto.Message
	Handler func(ctx context.Context, in proto.Message) (proto.Message, error)
}
