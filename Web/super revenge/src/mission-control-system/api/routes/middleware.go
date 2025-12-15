package routes

import (
	"app-api/api/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// ApplyAuthMiddleware applies the authentication middleware to a router
func ApplyAuthMiddleware(router *mux.Router, logger *zap.Logger) *mux.Router {
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(logger))
	return protectedRouter
}

func ApplyAdminMiddleware(router *mux.Router, logger *zap.Logger) *mux.Router {
	adminRouter := router.NewRoute().Subrouter()
	adminRouter.Use(middleware.AdminMiddleware(logger))
	return adminRouter
}
