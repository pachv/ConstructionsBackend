package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}
