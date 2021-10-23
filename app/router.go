package app

import "github.com/gin-gonic/gin"

func (app *App) Router() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.Use(app.Authenticate())
	{
		v1.GET("/containers/image/:containerName", app.GetContainerImage)

		v1.PUT("/images/pull", app.PullImage)
		v1.PUT("/containers/update", app.UpdateContainer)
		v1.PUT("/containers/rollback", app.RollbackContainer)
	}

	return router
}
