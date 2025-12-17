package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadResponse(c *gin.Context, errMessage string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": "error",
		"error":  errMessage,
	})
}
