package rest

import (
	"minecraft_searcher/api/rest/controllers"
	"minecraft_searcher/api/rest/middleware"
	"minecraft_searcher/lossesring"
	"minecraft_searcher/scanner"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(telemetry scanner.Telemetry, latestErrors lossesring.LossesRing[string]) {
	server := gin.New()

	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	api := server.Group("/api")

	api.POST("/login", controllers.Login)
	api.POST("/refresh", controllers.Refresh)

	protectedApi := server.Group("/api")
	protectedApi.Use(middleware.JwtAuthMiddleware)

	protectedApi.GET("/errors", func(c *gin.Context) {
		c.JSON(http.StatusOK, latestErrors.GetArray())
	})
	protectedApi.GET("/counter/:name", func(c *gin.Context) {
		switch c.Param("name") {
		case "workerInput":
			c.JSON(http.StatusOK, telemetry.GetWorkerInputDataSpeed())
		default:
			c.Status(http.StatusNotFound)
		}
	})
	protectedApi.GET("/scanners", func(c *gin.Context) {
		workers := make(chan []scanner.Worker)
		telemetry.GetController().GetWorkers <- workers
		c.JSON(http.StatusOK, <-workers)
	})

	server.Run()
}
