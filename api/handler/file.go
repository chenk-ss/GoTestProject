package handler

import (
	"goTestProject/tools"

	"github.com/gin-gonic/gin"
)

type QueryFileInfoParam struct {
	Id string `form:"id" json:"id" binding:"required"`
}

func FileInfo(c *gin.Context) {
	var param QueryFileInfoParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	tools.SuccessWithMsg(c, "Query file info ok!", param.Id)
}
