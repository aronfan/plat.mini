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

func BadRequestV2(c *gin.Context, data any) {
	if data != nil {
		FailResp(c, http.StatusBadRequest, data)
	} else {
		FailResp(c, http.StatusBadRequest, "Bad Request")
	}
}

func InternalError(c *gin.Context) {
	FailResp(c, http.StatusInternalServerError, "Internal Server Error")
}

func InternalErrorV2(c *gin.Context, data any) {
	if data != nil {
		FailResp(c, http.StatusInternalServerError, data)
	} else {
		FailResp(c, http.StatusInternalServerError, "Internal Server Error")
	}
}

func InProgress(c *gin.Context) {
	FailResp(c, http.StatusAccepted, "Data processing in progress")
}
