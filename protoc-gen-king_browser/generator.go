package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/juju/errors"
)

type generator struct{}

func (g *generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resp := new(plugin.CodeGeneratorResponse)

	for _, name := range req.FileToGenerate {
		file, err := getFileDescriptor(req, name)
		if err != nil {
			return nil, err
		}

		genFile, err := g.generateFile(file)
		if err != nil {
			return nil, errors.Annotatef(err, "generating %s", name)
		}

		resp.File = append(resp.File, genFile)
	}

	return resp, nil
}

func (g *generator) generateFile(file *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	buffer := new(bytes.Buffer)

	tmpl, err := template.New("browser_file").Parse(browserTemplate)
	if err != nil {
		return nil, errors.Trace(err)
	}

	data := &templateData{
		proto:          file,
		Version:        Version,
		SourceFilename: file.GetName(),
		Package:        file.GetPackage(),
	}
	if err = tmpl.Execute(buffer, data); err != nil {
		return nil, errors.Annotatef(err, "rendering template for %s", file.GetName())
	}

	name := file.GetName()
	name = name[:len(name)-len(filepath.Ext(name))]
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fmt.Sprintf("%s.king.js", name)),
		Content: proto.String(buffer.String()),
	}, nil
}

func getFileDescriptor(req *plugin.CodeGeneratorRequest, name string) (*descriptor.FileDescriptorProto, error) {
	for _, descriptor := range req.ProtoFile {
		if descriptor.GetName() == name {
			return descriptor, nil
		}
	}

	return nil, errors.Errorf("could not find descriptor for %s", name)
}
