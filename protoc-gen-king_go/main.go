package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/juju/errors"
)

// Version is the release version of this generator
const Version = "1.0.0"

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		log.Println(Version)
		return
	}

	req, err := readCodeGeneratorRequest(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if len(req.FileToGenerate) == 0 {
		log.Fatal("no files to generate")
	}

	gen := new(generator)
	resp, err := gen.Generate(req)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeCodeGeneratorResponse(os.Stdout, resp); err != nil {
		log.Fatal(err)
	}
}

func readCodeGeneratorRequest(in io.Reader) (*plugin.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, errors.Trace(err)
	}

	req := new(plugin.CodeGeneratorRequest)
	if err := proto.Unmarshal(data, req); err != nil {
		return nil, errors.Trace(err)
	}

	return req, nil
}

func writeCodeGeneratorResponse(out io.Writer, resp *plugin.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return errors.Trace(err)
	}

	if _, err := out.Write(data); err != nil {
		return errors.Trace(err)
	}

	return nil
}
