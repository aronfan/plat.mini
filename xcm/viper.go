package xcm

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func NewContext() *viper.Viper {
	return viper.New()
}

func LoadConfigFileWithContext(filename string, ctx *viper.Viper) error {
	if ctx == nil {
		ctx = viper.GetViper()
	}

	ctx.SetConfigFile(filename)
	if err := ctx.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func LoadConfigFile(filename string) error {
	return LoadConfigFileWithContext(filename, nil)
}

func MapToStructWithContext[T any](ctx *viper.Viper) (*T, error) {
	if ctx == nil {
		ctx = viper.GetViper()
	}

	t := new(T)
	if err := ctx.Unmarshal(t); err != nil {
		return nil, err
	}
	return t, nil
}

func MapToStruct[T any]() (*T, error) {
	return MapToStructWithContext[T](nil)
}

func BeginWatchConfigWithContext(cb func(), ctx *viper.Viper) {
	if ctx == nil {
		ctx = viper.GetViper()
	}

	merger := NewTimeoutBasedMerger[fsnotify.Event]()
	onEvent := merger.Start(func([]*fsnotify.Event) { cb() })

	ctx.OnConfigChange(onEvent)
	ctx.WatchConfig()
}

func BeginWatchConfig(cb func()) {
	BeginWatchConfigWithContext(cb, nil)
}
