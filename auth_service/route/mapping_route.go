package routes

import "github.com/julienschmidt/httprouter"

type Route struct {
	Method string
	Path   string
}

type Router struct {
	Router *httprouter.Router
	Routes []Route
}

func New() *Router {
	return &Router{
		Router: httprouter.New(),
	}
}

func (route *Router) add(method, path string, handler httprouter.Handle) {
	route.Routes = append(route.Routes, Route{
		Method: method,
		Path:   path,
	})

	route.Router.Handle(method, path, handler)
}

func (route *Router) GET(path string, handle httprouter.Handle) {
	route.add("GET", path, handle)
}

func (route *Router) POST(path string, handle httprouter.Handle) {
	route.add("POST", path, handle)
}

func (route *Router) PUT(path string, handle httprouter.Handle) {
	route.add("PUT", path, handle)
}

func (route *Router) DELETE(path string, handle httprouter.Handle) {
	route.add("DELETE", path, handle)
}
