package response

import "github.com/gin-gonic/gin"

func Success(ctx *gin.Context, msg string, data interface{}){
	ctx.JSON(200, gin.H{
		"success": true,
		"message": msg,
		"data": data,
	})
}

func Failed(ctx *gin.Context, code int, msg *string, data ...interface{}){
	var d interface{} = nil
	var m string = "Something when wrong."

	if len(data) > 0 {
		d = data[0]
	}
	if msg != nil {
		m = *msg
	}

	ctx.JSON(code, gin.H{
		"success": false,
		"message": m,
		"data": d,
	})
}