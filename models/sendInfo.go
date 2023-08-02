package models

type SendInfo struct {
	Id         uint     `json:"-" uri:"msgId" gorm:"primary_key"`
	Msg        string   `json:"msg" gorm:"not null"`      // 要发送的消息
	Link       string   `json:"link"`                     // 点击跳转链接
	Users      []string `json:"users" gorm:"-"`           // 发送给谁，使用字符串数组存储
	UsersDB    string   `json:"usersDB" gorm:"type:text"` // 发送给谁，使用字符串数组存储
	Sender     string   `json:"sender" gorm:"not null"`   // 发送者
	Platform   string   `json:"platform" gorm:"not null"` // 发送到哪个平台
	Tenant     *Tenant  `json:"-" gorm:"foreignKey:TenantId;references:Id"`
	TenantId   uint     `json:"-" gorm:"index"`          // 租户ID，使用索引
	MsgType    string   `json:"msgType" gorm:"not null"` // 消息类型：图文，markdown，文字
	Title      string   `json:"title" gorm:"not null"`   // 标题
	TitleColor string   `json:"titleColor" `             // 标题颜色
	PngLink    string   `json:"pngLink"`                 // 消息图片链接
	IsRead     bool     `json:"-" gorm:"default:false"`
}

func (SendInfo) TableName() string {
	return "message"
}
