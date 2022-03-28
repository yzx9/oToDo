package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/entity"
)

func SendSmsCode(c *gin.Context) {
	phone, ok := c.GetQuery("phone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "手机号不存在",
		})
		return
	}
	isSend := bll.SendCode(phone)
	if isSend {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "发送成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "发送成功",
		})
	}

}
func SmsLogin(c *gin.Context) {
	var smsparam entity.SmsLoginParam
	err := c.BindJSON(&smsparam)
	if err != nil {
		panic(err)
	}
	user := bll.SmsLogin(smsparam)
	if user != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "登陆成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "登录失败",
	})

}
