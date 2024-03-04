package xtest

type testConfig struct {
	HttpConfig struct {
		Addr  string `viper:"addr"`
		Debug bool   `viper:"debug"`
	} `mapstructure:"http"`
	LogConfig struct {
		Path string `viper:"path"`
	} `mapstructure:"log"`
	RedisConfig struct {
		Url string `viper:"url"`
		Ssh struct {
			Enable  bool   `viper:"enable"`
			Addr    string `viper:"addr"`
			User    string `viper:"user"`
			Keyfile string `viper:"keyfile"`
			Keypass string `viper:"keypass"`
		} `mapstructure:"ssh"`
	} `mapstructure:"redis"`
}
