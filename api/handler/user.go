package handler

import (
	"encoding/json"
	"goTestProject/config"
	"goTestProject/db"
	"goTestProject/logic/dao"
	"goTestProject/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var userDao = dao.User{}
var redisClient = db.RedisClient

type UserDetail struct {
	Id   int64  `json:"id"`
	Name string `json:"username"`
	Age  int    `json:"age"`
}

type QueryUserInfoParam struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func UserInfo(c *gin.Context) {
	var param QueryUserInfoParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	user := userDao.QueryUserByName(param.Name)
	if user.Id == 0 {
		tools.FailWithMsg(c, "Query user error! User not existing!")
		return
	}
	userDetail := UserDetail{
		user.Id, user.Name, time.Now().Local().Year() - user.Birthday.Local().Year(),
	}
	tools.SuccessWithMsg(c, "Query user info ok!", userDetail)
}

type UserRegisterParam struct {
	Name     string    `form:"name" json:"name" binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday" binding:"required"`
}

func Register(c *gin.Context) {
	var param UserRegisterParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	user, err := userDao.Add(param.Name, param.Birthday)
	if err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	tools.SuccessWithMsg(c, "Add user success!", user.Id)
	if bts, err := json.Marshal(user); err == nil {
		redisClient.HSet(config.RedisPrefix+"user", strconv.FormatInt(user.Id, 10), string(bts))
	}
}
