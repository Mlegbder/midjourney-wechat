package model

import (
	"time"
)

type WUser struct {
	Id        int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"id"`
	WechatId  string    `gorm:"column:wechat_id;type:varchar(255);comment:微信ID;NOT NULL" json:"wechat_id"`
	NickName  string    `gorm:"column:nick_name;type:varchar(255);comment:昵称" json:"nick_name"`
	Balance   int       `gorm:"column:balance;type:int(11);default:0;comment:余额;NOT NULL" json:"balance"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;NOT NULL" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt"`
	IsAdmin   string    `gorm:"column:is_admin;type:char;comment:是否管理员" json:"isAdmin"`
}

func (m *WUser) TableName() string {
	return "w_user"
}
