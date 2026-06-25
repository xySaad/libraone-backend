package trail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Router[T any] struct {
	mux        *http.ServeMux
	middleware Middleware[T]
}

// ServeHTTP implements [http.Handler].
func (rtr Router[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.mux.ServeHTTP(w, r)
}

type Middleware[T any] func(c *Context) (T, *Error)
type Handler[T any] func(c *Context, dep T) (Success, *Error)

func NewRouter[T any](middleware Middleware[T]) Router[T] {
	return Router[T]{mux: http.NewServeMux(), middleware: middleware}
}

type Null struct{}
type NoDepHandler = func(c *Context, _ Null) (Success, *Error)

func NoOpMiddleware() Middleware[Null] {
	return func(c *Context) (Null, *Error) { return Null{}, nil }
}

func DefaultRouter() Router[Null] {
	return Router[Null]{
		mux:        http.NewServeMux(),
		middleware: NoOpMiddleware(),
	}
}

func (rtr *Router[T]) AddRoute(pattern string, handler Handler[T]) {
	rtr.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{Request: r, respWr: w, context: context.Background()}
		reqHeader := fmt.Sprintf("%s %s\n", r.Method, r.RequestURI)

		middlewareResult, err := rtr.middleware(ctx)
		if err != nil {
			w.WriteHeader(err.Status)
			fmt.Fprintf(w, "%d - %s\n", err.Status, err.Message)
			fmt.Fprintf(os.Stderr, "%s\n[%d] %s\n%s\n", reqHeader, err.Status, err.Message, err.error)
			return
		}

		handlerResult, err := handler(ctx, middlewareResult)
		if err != nil {
			w.WriteHeader(err.Status)
			fmt.Fprintf(w, "%d - %s\n", err.Status, err.Message)
			fmt.Fprintf(os.Stderr, "%s\n[%d] %s\n%s\n", reqHeader, err.Status, err.Message, err.error)
			return
		}

		for key, values := range handlerResult.Headers {
			w.Header()[key] = values
		}
		w.WriteHeader(handlerResult.Status)

		if bodyReader, ok := handlerResult.Body.(io.ReadCloser); ok {
			defer bodyReader.Close()
			_, err := io.Copy(w, bodyReader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n[%s] %s\n", reqHeader, "FATAL", err)
			}

			fmt.Fprintf(os.Stdout, "[%d] %s\n", handlerResult.Status, reqHeader)
			return
		}

		err2 := json.NewEncoder(w).Encode(handlerResult.Body)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%d - %s\n", http.StatusInternalServerError, "error writing response")
			fmt.Fprintf(os.Stderr, "%s\n[%d] %s\n", reqHeader, http.StatusInternalServerError, err2)
			return
		}

		fmt.Fprintf(os.Stdout, "[%d] %s\n", handlerResult.Status, reqHeader)
	})
}

// Extend creates a new router that wraps the passed middleware along with the parent router's middleware
func Extend[R any, T any](parent Router[R], middleware Middleware[T]) Router[T] {
	return Router[T]{mux: parent.mux, middleware: func(c *Context) (T, *Error) {
		//TODO(xySaad): pass down the parent middleware result on sucess
		_, err := parent.middleware(c)
		if err != nil {
			var result T
			return result, err
		}
		return middleware(c)
	}}
}

var _ http.Handler = (*Router[any])(nil)
