package king

import (
	"context"
	"net/http"

	"github.com/altipla-consulting/sentry"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"

	"github.com/altipla-consulting/king/peer"
	"github.com/altipla-consulting/king/runtime"
)

func WithHttprouter(router *httprouter.Router) ServerOption {
	return func(server *Server) {
		server.router = router
	}
}

func Debug(debug bool) ServerOption {
	return func(server *Server) {
		server.debug = debug
	}
}

func WithLogrus() ServerOption {
	return func(server *Server) {
		server.errorMiddlewares = append(server.errorMiddlewares, func(ctx context.Context, appErr error) {
			if server.debug {
				log.WithFields(log.Fields{"err": appErr.Error()}).Error("call failed")
				log.Error("Error stack:\n", errors.ErrorStack(appErr))
			} else {
				log.WithFields(log.Fields{"err": appErr.Error(), "stack": errors.ErrorStack(appErr)}).Error("call failed")
			}
		})
	}
}

func WithSentry(dsn string) ServerOption {
	client := sentry.NewClient(dsn)

	return func(server *Server) {
		server.errorMiddlewares = append(server.errorMiddlewares, func(ctx context.Context, appError error) {
			r := peer.RequestFromContext(ctx)

			client.ReportRequest(appError, r)
		})
	}
}

func WithHttpClient(client *http.Client) runtime.ClientOption {
	return func(caller *runtime.ClientCaller) {
		caller.Client = client
	}
}

func WithAuthorization(token string) runtime.ClientOption {
	return func(caller *runtime.ClientCaller) {
		caller.Authorization = token
	}
}

func WithServerTraceOption(traceOption trace.StartOption) ServerOption {
	return func(server *Server) {
		server.traceOptions = append(server.traceOptions, traceOption)
	}
}

func WithClientTraceOption(traceOption trace.StartOption) runtime.ClientOption {
	return func(caller *runtime.ClientCaller) {
		caller.TraceOptions = append(caller.TraceOptions, traceOption)
	}
}
