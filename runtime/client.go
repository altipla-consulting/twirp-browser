package runtime

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

type ClientCaller struct {
	Server string
	Client *http.Client
}

type ClientOption func(clientCaller *ClientCaller)

func (caller *ClientCaller) Call(ctx context.Context, serviceName, methodName string, in, out proto.Message) error {
	serialized, err := proto.Marshal(in)
	if err != nil {
		return errors.Trace(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/_/%s/%s", caller.Server, serviceName, methodName), bytes.NewReader(serialized))
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Add("Content-Type", "application/protobuf")
	req.Header.Add("Accept", "application/protobuf")

	hc := http.DefaultClient
	if caller.Client != nil {
		hc = caller.Client
	}
	resp, err := hc.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Annotatef(errors.New("unexpected status code"), "status: %s", resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

	if err := proto.Unmarshal(content, out); err != nil {
		return errors.Trace(err)
	}

	return nil
}
