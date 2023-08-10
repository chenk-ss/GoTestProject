package handler

import (
	"goTestProject/tools"

	"github.com/gin-gonic/gin"
)

type CreateKafkaMsgParam struct {
	Msg string `form:"msg" json:"msg" binding:"required"`
}

func PushMsg(c *gin.Context) {
	var param CreateKafkaMsgParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	tools.PushMsgToKafka(param.Msg)
	tools.SuccessWithMsg(c, "Push msg success!", param.Msg)
}
