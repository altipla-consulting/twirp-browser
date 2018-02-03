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

func ClientCall(ctx context.Context, server, serviceName, methodName string, in, out proto.Message) error {
	serialized, err := proto.Marshal(in)
	if err != nil {
		return errors.Trace(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/_/%s/%s", server, serviceName, methodName), bytes.NewReader(serialized))
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Add("Content-Type", "application/protobuf")
	req.Header.Add("Accept", "application/protobuf")

	resp, err := http.DefaultClient.Do(req)
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
