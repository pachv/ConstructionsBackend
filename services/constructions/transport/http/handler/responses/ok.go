package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OkResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})

}
