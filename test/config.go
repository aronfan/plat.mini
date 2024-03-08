package test

type sshConfig struct {
	Enable  bool   `viper:"enable"`
	Addr    string `viper:"addr"`
	User    string `viper:"user"`
	Keyfile string `viper:"keyfile"`
	Keypass string `viper:"keypass"`
}

type testConfig struct {
	LogConfig struct {
		Path  string `viper:"path"`
		Level string `viper:"level"`
	} `mapstructure:"log"`
	HttpConfig struct {
		Addr  string `viper:"addr"`
		Debug bool   `viper:"debug"`
	} `mapstructure:"http"`
	MysqlConfig struct {
		Dsn string    `viper:"dsn"`
		Ssh sshConfig `mapstructure:"ssh"`
	} `mapstructure:"mysql"`
	RedisConfig struct {
		Url string    `viper:"url"`
		Ssh sshConfig `mapstructure:"ssh"`
	} `mapstructure:"redis"`
}
