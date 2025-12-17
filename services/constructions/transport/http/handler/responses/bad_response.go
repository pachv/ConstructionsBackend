package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalServiceErrorResponse(c *gin.Context, err string) {

	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  err,
	})

}
