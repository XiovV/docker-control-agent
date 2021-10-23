package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) successResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (app *App) badRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"error": message})
}

func (app *App) notFoundErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func (app *App) internalErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}
