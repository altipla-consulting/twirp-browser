package runtime

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

type ClientCaller struct {
	Server        string
	Client        *http.Client
	Authorization string
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", caller.Authorization))

	var duration time.Duration
	if dl, ok := ctx.Deadline(); ok {
		duration := dl.Sub(time.Now())
		if duration > 25*time.Second {
			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*25))
			defer cancel()

			dl, _ = ctx.Deadline()
		}

		dlout, err := dl.MarshalText()
		if err != nil {
			return errors.Trace(err)
		}
		req.Header.Add("X-King-Deadline", fmt.Sprintf("%v", dlout))
	}

	hc := &http.Client{
		Timeout: duration,
	}
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
