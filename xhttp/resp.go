package xhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccResp(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func FailResp(c *gin.Context, code int, data any) {
	c.JSON(code, data)
}

func BadRequest(c *gin.Context) {
	FailResp(c, http.StatusBadRequest, "Bad Request")
}

func InternalError(c *gin.Context) {
	FailResp(c, http.StatusInternalServerError, "Internal Server Error")
}

func InProgress(c *gin.Context) {
	FailResp(c, http.StatusAccepted, "Data processing in progress")
}
