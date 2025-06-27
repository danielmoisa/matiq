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
	flowActionRouter := routerGroup.Group("/teams/:teamID/workflow/:workflowID/flowActions")

	// flowActionRouter.Use(remotejwtauth.RemoteJWTAuth())

	// health router
	healthRouter.GET("", r.Controller.GetHealth)

	// flow action routers
	flowActionRouter.POST("", r.Controller.CreateWorkflow)
	// flowActionRouter.GET("/:flowActionID", r.Controller.GetFlowAction)
	// flowActionRouter.PUT("/:flowActionID", r.Controller.UpdateFlowAction)
	// flowActionRouter.DELETE("/:flowActionID", r.Controller.DeleteFlowAction)
	// flowActionRouter.POST("/:flowActionID/run", r.Controller.RunFlowAction)
	// flowActionRouter.PUT("/byBatch", r.Controller.UpdateFlowActionByBatch)
}
