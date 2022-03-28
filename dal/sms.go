package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertSmsCode(smscode *entity.SmsCode) error {
	err := db.Create(smscode)
	return util.WrapGormErr(err.Error, "user")
}
func SelecSmsCode(id int64) (entity.SmsCode, error) {
	var smscode entity.SmsCode
	err := db.Where("ID=?", id).First(&smscode).Error
	return smscode, util.WrapGormErr(err, "smscode")

}
func UpdateCode(id int64, code string) error {
	var smscode entity.SmsCode
	err := db.Where("ID=?", id).First(&smscode).Error
	if err != nil {
		return util.WrapGormErr(err, "smscode")
	}
	err = db.Model(&smscode).Updates(map[string]interface{}{"Code": code}).Error
	return util.WrapGormErr(err, "smscode")
}
func DeleteSmsCode(id int64) error {
	err := db.Where("ID=?", id).Delete(entity.SmsCode{}).Error
	return util.WrapGormErr(err, "smscode")
}
func ValidateSmsCode(phone string, code string) *entity.SmsCode {
	var sms entity.SmsCode
	if err := db.Where("phone=? and code=? ", phone, code).First(&sms).Error; err != nil {
		panic(err)
	} else {
		return &sms
	}

}
func QueryByPhone(phone string) entity.User {
	var user entity.User
	if err := db.Where("Telephone=?", phone).First(&user).Error; err != nil {
		panic(err)
	}
	return user
}
