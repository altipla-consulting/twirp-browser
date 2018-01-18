package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type templateData struct {
	proto *descriptor.FileDescriptorProto

	Version        string
	Package        string
	SourceFilename string
}

func (p *templateData) Services() []*templateService {
	services := make([]*templateService, len(p.proto.GetService()))

	for i, svc := range p.proto.GetService() {
		services[i] = &templateService{
			proto: svc,
			pkg:   p.proto.GetPackage(),
		}
	}

	return services
}

func (p *templateData) Quote() string {
	return "`"
}

type templateMessage struct {
	proto *descriptor.DescriptorProto
	pkg   string
}

type templateService struct {
	proto *descriptor.ServiceDescriptorProto
	pkg   string
}

func (svc *templateService) ClientName() string {
	return svc.proto.GetName() + "Client"
}

func (svc *templateService) ServiceName() string {
	if svc.pkg == "" {
		return svc.proto.GetName()
	}

	return svc.pkg + "." + svc.proto.GetName()
}

func (svc *templateService) Methods() []*templateMethod {
	methods := make([]*templateMethod, len(svc.proto.GetMethod()))

	for i, proto := range svc.proto.GetMethod() {
		methods[i] = &templateMethod{
			proto: proto,
		}
	}

	return methods
}

type templateMethod struct {
	proto *descriptor.MethodDescriptorProto
}

func (m *templateMethod) MethodName() string {
	return m.proto.GetName()
}

const browserTemplate = `// Code generated by protoc-gen-king_browser {{.Version}}, DO NOT EDIT.
// Source: {{.SourceFilename}}
{{range .Services}}
class {{.ClientName}} {
	constructor({server = '', authorization = '', hook = null} = {}) {
		this.server = server;
		this.authorization = authorization;
		this.hook = hook;
	}
	{{range .Methods}}
	{{.MethodName}}(req) {
		return this._doRequest('{{.MethodName}}', req);
	}
	{{end}}
	_doRequest(method, req) {
		let opts = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(req),
		};
		if (this.authorization) {
			opts.headers.Authorization = this.authorization;
		}
		return fetch({{$.Quote}}${this.server}/_/{{.ServiceName}}/${method}{{$.Quote}}, opts)
			.then(resp => {
				if (resp.status !== 200) {
					let err = new Error(response.statusText);
					err.response = response;
					throw err;
				}

				if (hook) {
					hook(response);
				}

				return response.json();
			});
	}
}
{{end}}`