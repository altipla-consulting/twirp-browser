package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
	"golang.org/x/net/context"

	"github.com/altipla-consulting/king/internal/httperr"
)

type ClientCaller struct {
	Server        string
	Client        *http.Client
	Authorization string
	TraceOptions  []trace.StartOption
}

type ClientOption func(clientCaller *ClientCaller)

func NewClientCaller(server string) *ClientCaller {
	return &ClientCaller{
		Server: server,
		TraceOptions: []trace.StartOption{
			trace.WithSpanKind(trace.SpanKindClient),
		},
	}
}

func (caller *ClientCaller) Call(ctx context.Context, serviceName, methodName string, in, out proto.Message) error {
	ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s.%s", serviceName, methodName), caller.TraceOptions...)
	defer span.End()

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

	f := new(propagation.HTTPFormat)
	f.SpanContextToRequest(span.SpanContext(), req)

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
		req.Header.Add("X-King-Deadline", fmt.Sprintf("%s", dlout))
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

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

	if _, ok := httperr.StatusKingErr[resp.StatusCode]; ok {
		kingErr := new(KingError)
		if err := json.Unmarshal(content, kingErr); err != nil {
			return errors.Trace(err)
		}

		message := fmt.Sprintf("%s.%s", serviceName, methodName)

		switch resp.StatusCode {
		case http.StatusNotFound:
			return errors.NewNotFound(kingErr, message)
		case http.StatusUnauthorized:
			return errors.NewUnauthorized(kingErr, message)
		case http.StatusNotImplemented:
			return errors.NewNotImplemented(kingErr, message)
		case http.StatusBadRequest:
			return errors.NewBadRequest(kingErr, message)
		case http.StatusForbidden:
			return errors.NewForbidden(kingErr, message)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Annotatef(errors.New("unexpected status code"), "status: %s", resp.Status)
	}

	if err := proto.Unmarshal(content, out); err != nil {
		return errors.Trace(err)
	}

	return nil
}
