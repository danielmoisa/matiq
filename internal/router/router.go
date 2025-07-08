package router

import (
	"github.com/danielmoisa/workflow-builder/internal/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller *controller.Controller
}

func NewRouter(controller *controller.Controller) *Router {
	return &Router{
		Controller: controller,
	}
}

func (r *Router) RegisterRouters(engine *gin.Engine) {
	// config
	engine.UseRawPath = true

	// init route
	routerGroup := engine.Group("/api/v1")
	healthRouter := routerGroup.Group("/health")
	workflowRouter := routerGroup.Group("/workflow")
	authRouter := routerGroup.Group("/auth")

	// health router
	healthRouter.GET("", r.Controller.GetHealth)

	// Auth routes
	// Public -- no authentication required
	authRouter.POST("/login", r.Controller.Login)
	authRouter.POST("/register", r.Controller.Register)
	authRouter.POST("/refresh", r.Controller.RefreshToken)
	authRouter.POST("/logout", r.Controller.Logout)
	authRouter.GET("/validate", r.Controller.ValidateToken)

	// Auth routes
	// Protected -- requires Bearer token
	protectedAuthRouter := authRouter.Group("")
	protectedAuthRouter.Use(r.Controller.AuthMiddleware())
	protectedAuthRouter.GET("/profile", r.Controller.GetProfile)

	// Workflow routes
	// Protected -- requires Bearer token
	workflowRouter.Use(r.Controller.AuthMiddleware())
	workflowRouter.POST("", r.Controller.CreateWorkflow)
	workflowRouter.GET("/:workflowID", r.Controller.GetWorkflow)

}
