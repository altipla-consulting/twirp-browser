package runtime

import (
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
)

func CodecFromType(codecType string) Codec {
	switch codecType {
	case "application/json":
		return new(CodecJSON)
	case "application/protobuf":
		return new(CodecPB)
	}

	return new(CodecJSON)
}

type Codec interface {
	ContentType() string
	Decode(r io.Reader, data proto.Message) error
	Encode(w io.Writer, data proto.Message) error
}

type CodecJSON struct {
}

func (codec *CodecJSON) ContentType() string {
	return "application/json; charset=utf-8"
}

func (codec *CodecJSON) Decode(r io.Reader, data proto.Message) error {
	m := jsonpb.Unmarshaler{}
	return errors.Trace(m.Unmarshal(r, data))
}

func (codec *CodecJSON) Encode(w io.Writer, data proto.Message) error {
	m := jsonpb.Marshaler{
		EmitDefaults: true,
	}
	return errors.Trace(m.Marshal(w, data))
}

type CodecPB struct {
}

func (codec *CodecPB) ContentType() string {
	return "application/protobuf"
}

func (codec *CodecPB) Decode(r io.Reader, data proto.Message) error {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(proto.Unmarshal(content, data))
}

func (codec *CodecPB) Encode(w io.Writer, data proto.Message) error {
	content, err := proto.Marshal(data)
	if err != nil {
		return errors.Trace(err)
	}
	if _, err := w.Write(content); err != nil {
		return errors.Trace(err)
	}
	return nil
}
