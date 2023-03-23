package config

type AzureBlob struct {
	AccountName string `mapstructure:"account-name" json:"account-name" yaml:"account-name"`
	AccountKey  string `mapstructure:"account-key" json:"account-key" yaml:"account-key"`
	Container   string `mapstructure:"container" json:"container" yaml:"container"`
}
