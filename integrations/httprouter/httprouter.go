package httprouter

import (
	"fmt"
	"net/http"

	"github.com/altipla-consulting/king/peer"
	"github.com/altipla-consulting/king/runtime"
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
		inCodec := runtime.CodecFromType(r.Header.Get("Content-Type"))
		outCodec := runtime.CodecFromType(r.Header.Get("Accept"))

		w.Header().Set("Content-Type", outCodec.ContentType())

		r = r.WithContext(peer.RequestWithContext(r))

		in := method.Input()
		if err := inCodec.Decode(r.Body, in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out, err := method.Handler(r.Context(), in)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := outCodec.Encode(w, out); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}
}
