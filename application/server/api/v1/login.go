package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    "",
		"message": "string",
		"ok":      true,
	})
}
