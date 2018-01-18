package runtime

import (
	"github.com/golang/protobuf/proto"
)

type Hook func(data proto.Message) error
