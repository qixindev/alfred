package model

import "time"

type SendInfo struct {
	Id           uint      `json:"id" uri:"msgId" gorm:"primary_key"`
	Msg          string    `json:"msg" gorm:"not null"`      // 要发送的消息
	Link         string    `json:"link"`                     // 点击跳转链接
	Users        []string  `json:"users" gorm:"-"`           // 发送给谁，使用字符串数组存储
	UsersDB      string    `json:"usersDB" gorm:"type:text"` // 发送给谁，使用字符串数组存储
	ReceiverName string    `json:"receiverName" gorm:"-"`    // 接收者名字
	Sender       string    `json:"sender" gorm:"not null"`   // 发送者
	SenderName   string    `json:"senderName" gorm:"-"`      // 发送者名字
	Platform     string    `json:"platform" gorm:"not null"` // 发送到哪个平台
	Tenant       *Tenant   `json:"-" gorm:"foreignKey:TenantId;references:Id"`
	TenantId     uint      `json:"-" gorm:"index"`          // 租户ID，使用索引
	MsgType      string    `json:"msgType" gorm:"not null"` // 消息类型：图文，markdown，文字
	Title        string    `json:"title" gorm:"not null"`   // 标题
	TitleColor   string    `json:"titleColor" `             // 标题颜色
	PngLink      string    `json:"pngLink"`                 // 消息图片链接
	IsRead       bool      `json:"isRead" gorm:"default:false"`
	SendAt       time.Time `json:"sendAt" gorm:"type:timestamp default:CURRENT_TIMESTAMP"` // 发送时间
}

type SendInfoDB struct {
	Id           uint      `json:"id" uri:"msgId" gorm:"primary_key"`
	Msg          string    `json:"msg" gorm:"not null"`               // 要发送的消息
	Link         string    `json:"link"`                              // 点击跳转链接
	UsersDB      string    `json:"usersDB" gorm:"type:text"`          // 发送给谁，使用字符串数组存储
	ReceiverName string    `json:"receiverName" gorm:"receiver_name"` // 接收者名字
	Sender       string    `json:"sender" gorm:"not null"`            // 发送者
	SenderName   string    `json:"senderName" gorm:"sender_name"`     // 发送者名字
	Platform     string    `json:"platform" gorm:"not null"`          // 发送到哪个平台
	Tenant       *Tenant   `json:"-" gorm:"foreignKey:TenantId;references:Id"`
	TenantId     uint      `json:"-" gorm:"index"`          // 租户ID，使用索引
	MsgType      string    `json:"msgType" gorm:"not null"` // 消息类型：图文，markdown，文字
	Title        string    `json:"title" gorm:"not null"`   // 标题
	TitleColor   string    `json:"titleColor" `             // 标题颜色
	PngLink      string    `json:"pngLink"`                 // 消息图片链接
	IsRead       bool      `json:"isRead" gorm:"default:false"`
	SendAt       time.Time `json:"sendAt" gorm:"type:timestamp default:CURRENT_TIMESTAMP"` // 发送时间
}

func (SendInfo) TableName() string {
	return "message"
}
