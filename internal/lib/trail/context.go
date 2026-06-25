package trail

import (
	"context"
	"net/http"
	"time"
)

type Context struct {
	Request *http.Request
	respWr  http.ResponseWriter
	context context.Context
}

func (c *Context) Query(s string) string {
	return c.Request.URL.Query().Get(s)
}

// Deadline implements [context.Context].
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.context.Deadline()
}

// Done implements [context.Context].
func (c *Context) Done() <-chan struct{} {
	return c.context.Done()
}

// Err implements [context.Context].
func (c *Context) Err() error {
	return c.context.Err()
}

// Value implements [context.Context].
func (c *Context) Value(key any) any {
	return c.context.Value(key)
}

var _ context.Context = (*Context)(nil)

func (c *Context) Success(status int, headers http.Header, body any) (Success, *Error) {
	return Success{
		Status:  status,
		Headers: headers,
		Body:    body,
	}, nil
}

func (c *Context) Error(publicError PublicError, err error) (Success, *Error) {
	return Success{Status: publicError.Status},
		&Error{PublicError: publicError, error: err}
}
