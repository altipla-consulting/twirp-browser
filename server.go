package king

import (
	"fmt"
	"net/http"

	"github.com/altipla-consulting/collections"

	"github.com/altipla-consulting/king/peer"
	"github.com/altipla-consulting/king/runtime"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router           *httprouter.Router
	logrus           bool
	errorMiddlewares []ErrorMiddleware
	debug            bool
	cors             []string
}

type ErrorMiddleware func(appErr error)

type ServerOption func(server *Server)

func NewServer(opts ...ServerOption) *Server {
	server := new(Server)
	for _, opt := range opts {
		opt(server)
	}

	if server.router != nil {
		for _, svc := range runtime.Services {
			for _, method := range svc.Methods {
				path := fmt.Sprintf("/_/%s/%s", svc.Name, method.Name)
				server.router.POST(path, buildHandler(server, method))
				server.router.OPTIONS(path, corsHandler(server))
			}
		}
	}

	return server
}

func buildHandler(server *Server, method *runtime.Method) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		inCodec := runtime.CodecFromType(r.Header.Get("Content-Type"))
		outCodec := runtime.CodecFromType(r.Header.Get("Accept"))

		w.Header().Set("Content-Type", outCodec.ContentType())

		r = r.WithContext(peer.RequestWithContext(r))

		in := method.Input()
		if err := inCodec.Decode(r.Body, in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			for _, m := range server.errorMiddlewares {
				m(err)
			}
			return
		}

		out, err := method.Handler(r.Context(), in)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			for _, m := range server.errorMiddlewares {
				m(err)
			}
			return
		}

		if err := outCodec.Encode(w, out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			for _, m := range server.errorMiddlewares {
				m(err)
			}
			return
		}

		return
	}
}

func corsHandler(server *Server) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		origin := r.Header.Get("Origin")
		if collections.HasString(server.cors, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
