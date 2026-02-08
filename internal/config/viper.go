package config

import "github.com/spf13/viper"

func NewViper() (*viper.Viper, error) {
	v := viper.New()
	
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	v.AutomaticEnv()

	err := v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return v, nil
}