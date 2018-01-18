package httprouter

import (
	"fmt"
	"net/http"

	"github.com/altipla-consulting/king/peer"
	"github.com/altipla-consulting/king/runtime"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
)

func RegisterServices(r *httprouter.Router) {
	for _, svc := range runtime.Services {
		for _, method := range svc.Methods {
			r.POST(fmt.Sprintf("/_/%s/%s", svc.Name, method.Name), buildHandler(method))
		}
	}
}

func buildHandler(method *runtime.Method) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := safeHandler(method, w, r); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func safeHandler(method *runtime.Method, w http.ResponseWriter, r *http.Request) error {
	inCodec := runtime.CodecFromType(r.Header.Get("Content-Type"))
	outCodec := runtime.CodecFromType(r.Header.Get("Accept"))

	w.Header().Set("Content-Type", outCodec.ContentType())

	r = r.WithContext(peer.RequestWithContext(r))

	in := method.Input()
	if err := inCodec.Decode(r.Body, in); err != nil {
		return errors.Trace(err)
	}

	out, err := method.Handler(r.Context(), in)
	if err != nil {
		return errors.Trace(err)
	}

	if err := outCodec.Encode(w, out); err != nil {
		return errors.Trace(err)
	}

	return nil
}
