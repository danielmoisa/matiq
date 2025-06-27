package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
