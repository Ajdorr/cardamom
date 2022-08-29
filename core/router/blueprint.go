package router

import "github.com/gin-gonic/gin"

type Route struct {
	method  string
	handler gin.HandlerFunc
}

type Blueprint struct {
	path      string
	routes    map[string]Route
	subroutes []*Blueprint
}

func RegisterBlueprint(r *gin.Engine, bp *Blueprint) {

	for subpath, route := range bp.routes {
		r.Handle(route.method, bp.path+subpath, route.handler)
	}

	for _, subBp := range bp.subroutes {
		registerSubBlueprint(bp.path, r, subBp)
	}
}

func registerSubBlueprint(root string, r *gin.Engine, bp *Blueprint) {
	rootPath := root + bp.path

	for subpath, route := range bp.routes {
		r.Handle(route.method, rootPath+subpath, route.handler)
	}

	for _, subBp := range bp.subroutes {
		registerSubBlueprint(bp.path, r, subBp)
	}

}
