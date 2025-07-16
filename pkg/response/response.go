package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponseData struct {
	Code   int    `json:"code"`
	Err    string `json:"err"`
	Detail any    `json:"detail"`
}

func SuccessReponse(c *gin.Context, code int, data any) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msgSuccessMap[code],
		Data:    data,
	})
}

func ErrorReponse(c *gin.Context, code int, detail string) {
	c.JSON(code, ErrorResponseData{
		Code:   code,
		Err:    msgErrorMap[code],
		Detail: detail,
	})
}
