package runtime

import (
	"golang.org/x/net/context"
)

var Services []*Service

type Service struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name    string
	Handler func(ctx context.Context, inCodec, outCodec Codec, inHook, outHook Hook) error
}
