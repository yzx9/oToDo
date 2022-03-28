package entity

type SmsCode struct {
	ID         int64  `json:"id" gorm:"pk autoincr"`
	Phone      string `json:"phone" gorm:"varvhar(11)"`
	BizId      string `json:"biz_id" gorm:"varchar(30)"`
	Code       string `json:"code" gorm:"varchar(4)"`
	CreateTime int64  `json:"create_time" gorm:"bigint"`
}
type SmsLoginParam struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
