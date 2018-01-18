package runtime

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
)

type Codec interface {
	ContentType() string
	Decode(data proto.Message) error
	Encode(data proto.Message) error
}

type CodecJSON struct {
	w http.ResponseWriter
	r *http.Request
}

func (codec *CodecJSON) ContentType() string {
	return "application/json; charset=utf-8"
}

func (codec *CodecJSON) Decode(data proto.Message) error {
	m := jsonpb.Unmarshaler{}
	return errors.Trace(m.Unmarshal(codec.r.Body, data))
}

func (codec *CodecJSON) Encode(data proto.Message) error {
	m := jsonpb.Marshaler{
		EmitDefaults: true,
	}
	return errors.Trace(m.Marshal(codec.w, data))
}

type CodecPB struct {
	w http.ResponseWriter
	r *http.Request
}

func (codec *CodecPB) ContentType() string {
	return "application/protobuf"
}

func (codec *CodecPB) Decode(data proto.Message) error {
	content, err := ioutil.ReadAll(codec.r.Body)
	if err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(proto.Unmarshal(content, data))
}

func (codec *CodecPB) Encode(data proto.Message) error {
	content, err := proto.Marshal(data)
	if err != nil {
		return errors.Trace(err)
	}
	if _, err := codec.w.Write(content); err != nil {
		return errors.Trace(err)
	}
	return nil
}
