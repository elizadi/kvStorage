package router

import "github.com/gin-gonic/gin"

func Router(server *Server) *gin.Engine {
	router := gin.Default()
	handler := New(server.uc, server.log)
	api := router.Group("/key_value_storage")
	{
		api.GET("/:key", handler.GetValue)
		api.POST("/new_value", handler.CreateValue)
		api.POST("/value", handler.EditeValue)
		api.DELETE("/:key", handler.DeleteValue)
	}

	return router
}
