package xcm

type SshConfig struct {
	Enable  bool   `viper:"enable"`
	Addr    string `viper:"addr"`
	User    string `viper:"user"`
	Keyfile string `viper:"keyfile"`
	Keypass string `viper:"keypass"`
}

type AdminConfig struct {
	Addrs []string `viper:"addrs"`
}

type Config struct {
	Name      string `viper:"name"`
	Node      int64  `viper:"node"`
	LogConfig struct {
		Path  string `viper:"path"`
		Level string `viper:"level"`
	} `mapstructure:"log"`
	HttpConfig struct {
		Addr   string      `viper:"addr"`
		Static string      `viper:"static"`
		Debug  bool        `viper:"debug"`
		Admin  AdminConfig `mapstructure:"admin"`
	} `mapstructure:"http"`
	TokenConfig struct {
		Key     string `viper:"key"`
		Timeout int64  `viper:"timeout"`
	} `mapstructure:"token"`
	CacheConfig struct {
		Timeout int64 `viper:"timeout"`
	} `mapstructure:"cache"`
	MysqlConfig struct {
		Enable bool      `viper:"enable"`
		Dsn    string    `viper:"dsn"`
		Ssh    SshConfig `mapstructure:"ssh"`
	} `mapstructure:"mysql"`
	RedisConfig struct {
		Enable bool      `viper:"enable"`
		Url    string    `viper:"url"`
		Ssh    SshConfig `mapstructure:"ssh"`
	} `mapstructure:"redis"`
}
