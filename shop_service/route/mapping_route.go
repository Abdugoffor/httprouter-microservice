package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"shop_service/appctx"

	"github.com/julienschmidt/httprouter"
)

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

	// route pattern ni context ga yozib qo'yamiz
	// middleware shu context dan o'qiydi - params loop shart emas
	route.Router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), appctx.RoutePatternKey, path)
		handler(w, r.WithContext(ctx), ps)
	})
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

func (route *Router) SyncRoutes() {

	var data []map[string]string

	for _, r := range route.Routes {
		data = append(data, map[string]string{
			"service": "shop_service",
			"name":    r.Method + " " + r.Path,
		})
	}

	body := map[string]interface{}{
		"permissions": data,
	}

	jsonData, _ := json.Marshal(body)

	http.Post("http://localhost:8080/internal/permissions", "application/json", bytes.NewBuffer(jsonData))
}
