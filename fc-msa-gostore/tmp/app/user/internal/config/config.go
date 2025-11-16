package config

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Env  string `mapstructure:"env"`
	} `json:"app"`

	Logger struct {
		Level    string `mapstructure:"level"`    // Example: "info", "debug"
		Filepath string `mapstructure:"filepath"` // Example: "logs/account.log"
	} `mapstructure:"logger"`

	MongoDB struct {
		URI    string `mapstructure:"uri"`
		DbName string `mapstructure:"dbname"`
	} `mapstructure:"mongodb"`

	JWT struct {
		Secret string `mapstructure:"secret"`
		TTL    int    `mapstructure:"ttl"`
	} `mapstructure:"jwt"`
}
