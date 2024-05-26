package v1

import (
	"smallurl/internal/app/config"
	"smallurl/internal/app/delivery/http/middleware"
	"smallurl/pkg/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

const version = "v1"

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

type Routes []Route

func NewRouter(root string, routes Routes, mode config.Mode, l logger.Interface) (*gin.Engine, error) {
	if mode == config.Release || mode == config.ReleaseProf {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.RequestLogger(l), middleware.CheckPanic)
	rt := router.Group(root)
	v1 := rt.Group(version)

	for _, route := range routes {
		route.Middlewares = append(route.Middlewares, route.HandlerFunc)
		v1.Handle(route.Method, route.Pattern, route.Middlewares...)
	}

	if mode == config.DebugProf || mode == config.ReleaseProf {
		pprof.Register(router)
	}

	return router, nil
}
