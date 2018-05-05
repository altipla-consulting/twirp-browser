package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/altipla-consulting/collections"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/juju/errors"
)

type generator struct {
	imports map[string]string
	pkgs    map[string]string
}

func (g *generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resp := new(plugin.CodeGeneratorResponse)

	g.imports = make(map[string]string)
	g.pkgs = make(map[string]string)

	used := map[string]bool{}
	for _, file := range req.ProtoFile {
		pkg := filepath.Dir(file.GetName())
		if file.GetOptions().GetGoPackage() != "" {
			pkg = file.GetOptions().GetGoPackage()
		}

		name := filepath.Base(pkg)
		_, ok := used[name]
		for i := 2; ok; i++ {
			name = fmt.Sprintf("%s%d", filepath.Base(pkg), i)
			_, ok = used[name]
		}

		g.imports[file.GetPackage()] = fmt.Sprintf("%s %s", name, strconv.Quote(pkg))
		g.pkgs[file.GetPackage()] = name
	}

	for _, name := range req.FileToGenerate {
		file, err := getFileDescriptor(req, name)
		if err != nil {
			return nil, err
		}

		genFile, err := g.generateFile(file)
		if err != nil {
			return nil, errors.Annotatef(err, "generating %s", name)
		}

		if len(file.GetService()) > 0 {
			resp.File = append(resp.File, genFile)
		}
	}

	return resp, nil
}

func (g *generator) generateFile(file *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	buffer := new(bytes.Buffer)

	tmpl, err := template.New("go_file").Parse(browserTemplate)
	if err != nil {
		return nil, errors.Trace(err)
	}

	data := &templateData{
		proto:          file,
		Version:        Version,
		SourceFilename: file.GetName(),
		typesMap:       map[string]string{},
	}
	var methods bool
	for _, svc := range file.GetService() {
		for _, method := range svc.GetMethod() {
			methods = true
			inPkg, inName := splitType(method.GetInputType())
			if inPkg == file.GetPackage() {
				data.typesMap[method.GetInputType()] = inName
			} else {
				data.typesMap[method.GetInputType()] = fmt.Sprintf("%s.%s", g.pkgs[inPkg], inName)
				data.Imports = append(data.Imports, g.imports[inPkg])
			}

			outPkg, outName := splitType(method.GetOutputType())
			if outPkg == file.GetPackage() {
				data.typesMap[method.GetOutputType()] = outName
			} else {
				data.typesMap[method.GetOutputType()] = fmt.Sprintf("%s.%s", g.pkgs[outPkg], outName)
				data.Imports = append(data.Imports, g.imports[outPkg])
			}
		}
	}
	if methods {
		data.Imports = append(data.Imports, strconv.Quote("golang.org/x/net/context"))
		data.Imports = append(data.Imports, strconv.Quote("github.com/golang/protobuf/proto"))
		data.Imports = append(data.Imports, strconv.Quote("github.com/altipla-consulting/king/runtime"))
	}

	data.Package = strings.Replace(file.GetPackage(), ".", "_", -1)

	data.Imports = collections.UniqueStrings(data.Imports)
	sort.Strings(data.Imports)

	if err = tmpl.Execute(buffer, data); err != nil {
		return nil, errors.Annotatef(err, "rendering template for %s", file.GetName())
	}

	name := file.GetName()
	name = name[:len(name)-len(filepath.Ext(name))]
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fmt.Sprintf("%s.king.go", name)),
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

func splitType(typeName string) (string, string) {
	s := strings.Split(typeName, ".")
	return strings.Join(s[1:len(s)-1], "."), s[len(s)-1]
}
