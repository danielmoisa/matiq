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
	workflowRouter := routerGroup.Group("/teams/:teamID/workflow")

	// flowActionRouter.Use(remotejwtauth.RemoteJWTAuth())

	// health router
	healthRouter.GET("", r.Controller.GetHealth)

	// flow action routers
	workflowRouter.POST("", r.Controller.CreateWorkflow)
	workflowRouter.GET("/:workflowID", r.Controller.GetWorkflow)
	// flowActionRouter.PUT("/:flowActionID", r.Controller.UpdateFlowAction)
	// flowActionRouter.DELETE("/:flowActionID", r.Controller.DeleteFlowAction)
	// flowActionRouter.POST("/:flowActionID/run", r.Controller.RunFlowAction)
	// flowActionRouter.PUT("/byBatch", r.Controller.UpdateFlowActionByBatch)
}
