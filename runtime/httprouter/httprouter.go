package httprouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/altipla-consulting/king/runtime"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
)

type key int

const requestKey key = 0

func RegisterServices(r *httprouter.Router) {
	for _, svc := range runtime.Services {
		for _, method := range svc.Method {
			r.POST(fmt.Sprintf("/_/%s/%s", svc.Name, method.Name), buildHandler(method))
		}
	}
}

func buildHandler(method *runtime.Method, server interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var inCodec runtime.Codec
		switch r.Header.Get("Content-Type") {
		case "application/json":
			inCodec = &runtime.CodecJSON{w, r}
		case "application/protobuf":
			inCodec = &runtime.CodecPB{w, r}
		default:
			inCodec = &runtime.CodecJSON{w, r}
		}

		var outCodec runtime.Codec
		switch r.Header.Get("Accept") {
		case "application/json":
			outCodec = &runtime.CodecJSON{w, r}
		case "application/protobuf":
			outCodec = &runtime.CodecPB{w, r}
		default:
			outCodec = &runtime.CodecJSON{w, r}
		}

		w.Header().Set("Content-Type", outCodec.ContentType())

		r = r.WithContext(context.WithValue(ctx, requestKey, r))

		if err := method.Handler(ctx, inCodec, outCodec, inHook, outHook); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
