package bll

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"

	"github.com/yzx9/otodo/otodo"
)

func SendCode(phone string) bool {
	code := fmt.Sprintf("%0.6v", rand.New(rand.NewSource(time.Now().UnixMicro())).Int31n((1000000)))
	// smsConfig := otodo.ConfigSms{"cn-hangzhou", "LTAI5tR3J6TGs9ueDGNMejBo", "BFPxNvr6DZBgYvqiWEzVPSBx5U95fz", "阿里云短信测试", "SMS_154950909"}
	smsConfig, err := ParseConfig()
	if err != nil {
		panic(err)
	}
	client, err := dysmsapi.NewClientWithAccessKey(smsConfig.RegionID, smsConfig.AppKey, smsConfig.Appsecret)
	if err != nil {
		panic(err)
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = smsConfig.SignName
	request.TemplateCode = smsConfig.TemplateCode
	request.PhoneNumbers = phone
	par, err := json.Marshal(map[string]interface{}{
		"code": code,
	})
	request.TemplateParam = string(par)
	response, err := client.SendSms(request)
	fmt.Println(response)
	if err != nil {
		panic(err)
	}
	if response.Code == "OK" {
		//err:=DB.Save(),Hook
		smscode := entity.SmsCode{Code: code, Phone: phone, BizId: response.BizId, CreateTime: time.Now().Unix()}
		err := dal.InsertSmsCode(&smscode)
		if err != nil {
			panic(err)
		}
		return true

	} else {
		return false
	}

}
func SmsLogin(param entity.SmsLoginParam) *entity.User {
	sms := dal.ValidateSmsCode(param.Phone, param.Code)
	if sms == nil || time.Now().Unix()-sms.CreateTime > 300 {
		return nil
	}
	user := dal.QueryByPhone(param.Phone)
	return &user

}
func ParseConfig() (sms otodo.ConfigSms, err error) {
	var conf otodo.Config
	viper.SetConfigFile("./config.ymal")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file:%s \n", err))
	}
	if err = viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("Unmarshal conf failed err:%s \n", err))
	}
	viper.WatchConfig()
	sms = conf.Sms
	return sms, err

}
