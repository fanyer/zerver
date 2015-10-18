package route

import (
	"github.com/cosiner/zerver"
	"github.com/cosiner/zerver/utils/handle"
)

type Route struct {
	Pattern string

	Method       string
	HandleFunc   zerver.HandleFunc
	Interceptors []interface{}

	Handler interface{}
}

type Routes []Route

func convertHandleFunc(h interface{}) zerver.HandleFunc {
	switch h := h.(type) {
	case zerver.HandleFunc:
		return h
	case func(zerver.Request, zerver.Response):
		return h
	case func(zerver.Request, zerver.Response) error:
		return handle.Wrap(h)
	}

	panic("not a handle function")
}

func (r Routes) New(method, pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return append(r, Route{
		Method:       method,
		Pattern:      pattern,
		HandleFunc:   convertHandleFunc(handler),
		Interceptors: interceptors,
	})
}

func (r Routes) Apply(router zerver.Router) error {
	var err error
	for i := 0; err == nil && i < len(r); i++ {
		var route = r[i]
		if route.Handler != nil {
			err = router.Handle(route.Pattern, route.Handler)
		} else {
			err = router.HandleFunc(route.Pattern,
				route.Method,
				zerver.Intercept(route.HandleFunc,
					route.Interceptors...))
		}
	}

	return err
}

func (r Routes) Handle(pattern string, handler interface{}) Routes {
	return append(r, Route{
		Pattern: pattern,
		Handler: handler,
	})
}

func (r Routes) Get(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return r.New("GET", pattern, handler, interceptors...)
}

func (r Routes) Post(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return r.New("POST", pattern, handler, interceptors...)
}

func (r Routes) Delete(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return r.New("DELETE", pattern, handler, interceptors...)
}

func (r Routes) Put(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return r.New("PUT", pattern, handler, interceptors...)
}

func (r Routes) Patch(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return r.New("PATCH", pattern, handler, interceptors...)
}

func New(method, pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return Routes(nil).New(method, pattern, handler, interceptors...)
}

func Handle(pattern string, handler interface{}) Routes {
	return Routes(nil).Handle(pattern, handler)
}

func Get(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return New("GET", pattern, handler, interceptors...)
}

func Post(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return New("POST", pattern, handler, interceptors...)
}

func Delete(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return New("DELETE", pattern, handler, interceptors...)
}

func Put(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return New("PUT", pattern, handler, interceptors...)
}

func Patch(pattern string, handler interface{}, interceptors ...interface{}) Routes {
	return New("PATCH", pattern, handler, interceptors...)
}

func Apply(r zerver.Router, routes ...Route) error {
	return Routes(routes).Apply(r)
}
