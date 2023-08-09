package handler

import (
	"encoding/json"
	"fmt"
	"goTestProject/config"
	"goTestProject/db"
	"goTestProject/logic/dao"
	"goTestProject/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var userDao = dao.User{}
var redisClient = db.RedisClient

type UserDetail struct {
	Id       int64  `json:"id"`
	Name     string `json:"username"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type QueryUserInfoParam struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// Query user info
// @Summary      Query user info
// @Description  get user info by name
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name query string true "1"
// @Success      200  {object}  UserDetail
// @Router       /user/info [get]
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
		user.Id,
		user.Name,
		time.Now().Local().Year() - user.Birthday.Local().Year(),
		user.Password,
	}
	tools.SuccessWithMsg(c, "Query user info ok!", userDetail)
}

type UserRegisterParam struct {
	Name     string    `form:"name" json:"name" binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday" binding:"required"`
	Password string    `form:"password" json:"password" binding:"password"`
}

// Register a user
// @Summary      Register a user
// @Description  Register a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        param body UserRegisterParam true "params"
// @Success      200 string json "{"code":200,"data":{},"message":"success"}"
// @Router       /user/register [post]
func Register(c *gin.Context) {
	var param UserRegisterParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	user, err := userDao.Add(param.Name, param.Birthday, param.Password)
	if err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	tools.SuccessWithMsg(c, "Add user success!", user.Id)
	if bts, err := json.Marshal(user); err == nil {
		redisClient.HSet(config.RedisPrefix+"user", strconv.FormatInt(user.Id, 10), string(bts))
	}
}

type UserLoginParam struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var param UserLoginParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	user := userDao.QueryUserByName(param.Name)
	if user.Id == 0 {
		tools.FailWithMsg(c, "User not exist!")
		return
	}
	if user.Password != param.Password {
		tools.FailWithMsg(c, "User name or password not correct!")
		return
	}
	// idStr := strconv.FormatInt(user.Id, 10)
	sessionIndex := fmt.Sprintf("sess:userid:%d", user.Id)
	token := redisClient.Get(sessionIndex)
	if token.Val() != "" {
		tools.SuccessWithMsg(c, "Login success!", token.Val())
		return
	}
	randToken := tools.GetRandomToken(32)
	userData := make(map[string]interface{})
	userData["userId"] = user.Id
	userData["userName"] = user.Name

	redisClient.Do("MULTI")
	redisClient.HMSet("sess:session:"+randToken, userData)
	redisClient.Expire(randToken, 86400*time.Second)
	redisClient.Set(sessionIndex, randToken, 86400*time.Second)
	err := redisClient.Do("EXEC").Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return
	}
	redisClient.Set(randToken, userData, 86400*time.Second)
	tools.SuccessWithMsg(c, "Login success!", randToken)
}

func CheckToken(token string) bool {
	res := redisClient.Exists("sess:session:" + token)
	return res.Val() > 0
}
