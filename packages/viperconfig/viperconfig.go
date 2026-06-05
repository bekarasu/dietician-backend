package viperconfig

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const (
	configType      = "env"
	productionEnv   = "production"
	lookupStructTag = "mapstructure"
)

type Config struct {
	Env      string
	Path     string
	FileName string
}

func Load(c Config, i interface{}) (interface{}, error) {
	if c.Env == productionEnv {
		viper.AutomaticEnv()
		bindEnvs(viper.GetViper(), i)
	} else {
		viper.AddConfigPath(c.Path)
		viper.SetConfigName(c.FileName)
		viper.SetConfigType(configType)

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&i); err != nil {
		return nil, err
	}

	return i, nil
}

func bindEnvs(v *viper.Viper, ci interface{}) {
	iv := reflect.ValueOf(ci)
	it := reflect.TypeOf(ci)

	for i := 0; i < it.NumField(); i++ {
		fv := iv.Field(i)
		ft := it.Field(i)
		name := strings.ToLower(ft.Name)

		if tag, ok := ft.Tag.Lookup(lookupStructTag); ok {
			name = tag
		}

		if k := fv.Kind(); k == reflect.Struct {
			bindEnvs(v, fv.Interface())
		} else {
			_ = v.BindEnv(name)
		}
	}
}
