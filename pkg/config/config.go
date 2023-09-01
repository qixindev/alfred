package config

type Config struct {
	Zap        *Zap        `mapstructure:"zap" json:"zap" yaml:"zap"`
	Pgsql      *Pgsql      `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	System     *System     `mapstructure:"system" json:"system" yaml:"system"`
	AliyunOSS  *AliyunOSS  `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	TencentCOS *TencentCOS `mapstructure:"tencent-cos" json:"tencent-cos" yaml:"tencent-cos"`
	AzureBlob  *AzureBlob  `mapstructure:"azure-blob" json:"azure-blob" yaml:"azure-blob"`
	RabbitMq   *RabbitMq   `mapstructure:"rabbit-mq" json:"rabbit-mq" yaml:"rabbit-mq"`
	Urls       Urls        `mapstructure:"urls" json:"urls" yaml:"urls"`
}

type Urls struct {
	Airflow   string `mapstructure:"airflow" json:"airflow" yaml:"airflow"`
	AzureBlob string `mapstructure:"blob-url" json:"blob-url" yaml:"blob-url"`
}
