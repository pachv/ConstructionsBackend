package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestResponse(c *gin.Context, err string) {

	c.JSON(http.StatusBadRequest, gin.H{
		"status": "error",
		"error":  err,
	})

}
