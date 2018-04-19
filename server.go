package king

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/altipla-consulting/king/peer"
	"github.com/altipla-consulting/king/runtime"
	"github.com/altipla-consulting/sentry"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"

	"github.com/altipla-consulting/king/internal/httperr"
)

type Server struct {
	router           *httprouter.Router
	errorMiddlewares []ErrorMiddleware
	debug            bool
}

type ErrorMiddleware func(ctx context.Context, appErr error)

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
		r = r.WithContext(sentry.WithContext(r.Context()))

		in := method.Input()
		if err := inCodec.Decode(r.Body, in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			for _, m := range server.errorMiddlewares {
				m(r.Context(), err)
			}
			return
		}

		out, err := method.Handler(r.Context(), in)
		if err != nil {
			for _, m := range server.errorMiddlewares {
				m(r.Context(), err)
			}

			kingErr := &runtime.KingError{
				Message: err.Error(),
			}

			switch {
			case errors.IsNotFound(err):
				kingErr.Err = httperr.ErrorTypeNotFound
			case errors.IsUnauthorized(err):
				kingErr.Err = httperr.ErrorTypeUnauthorized
			case errors.IsNotImplemented(err):
				kingErr.Err = httperr.ErrorTypeNotImplemented
			case errors.IsBadRequest(err):
				kingErr.Err = httperr.ErrorTypeBadRequest
			case errors.IsForbidden(err):
				kingErr.Err = httperr.ErrorTypeForbidden
			default:
				kingErr.Err = httperr.ErrorTypeInternalServerError
			}

			data, err := json.Marshal(kingErr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				for _, m := range server.errorMiddlewares {
					m(r.Context(), err)
				}
				return
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, string(data), httperr.KingErrStatus[kingErr.Err])

			return
		}

		if err := outCodec.Encode(w, out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			for _, m := range server.errorMiddlewares {
				m(r.Context(), err)
			}
			return
		}

		return
	}
}
