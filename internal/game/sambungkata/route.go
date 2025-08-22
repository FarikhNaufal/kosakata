package sambungkata

import (

	"github.com/gin-gonic/gin"
)

func RegisterRoute(group *gin.RouterGroup, h *WordHandler) {
	group.GET("/", h.ShowAllWord).
		GET("/:id", h.ShowWord).
		POST("/store", h.StoreWord).
		GET("/today", h.GetTodayWord).
		POST("/check", h.CheckingWord)
}
