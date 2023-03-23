package config

type System struct {
	Env          string `mapstructure:"env" json:"env" yaml:"env"`                               // 环境
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`                            // 端口
	DbType       string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                   // 数据库类型:默认postgresql
	OssType      string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                // 对象存储类型
	UseRedis     bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`             // 使用redis
	LimitCountIP int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"` // ip限制数量
	LimitTimeIP  int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`    // ip限制时间
	Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`                            // 发布模式
}

func (s *System) IsDevMode() bool {
	return s.Mode == "dev"
}
