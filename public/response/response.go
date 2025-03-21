package response

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"requestId"`
	Data      interface{} `json:"data"`
}

func Success(context *gin.Context, data ...interface{}) {
	response := ApiResponse{}
	response.Code = 0
	response.Msg = "success"
	response.RequestId = context.GetString("requestId")

	if len(data) == 0 {
		response.Data = struct {
		}{}
	} else {
		response.Data = data[0]
	}
	context.JSON(http.StatusOK, response)
}

func SuccessStream(context *gin.Context, step func(w io.Writer) bool) {
	context.Header("Content-type", "application/octet-stream")
	context.Stream(step)
}

func SuccessEventStream(context *gin.Context, name string, data ...interface{}) {
	context.SSEvent(name, data[0])
}
