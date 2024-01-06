package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port 			string
	Sync_time 		time.Duration
	Num_processes 	int
	Url		 		string
	JwtSecretKey	string
	CookieName		string
	Postgres		Postgres
	Redis			Redis
	Cookie			Cookie
}

type Postgres struct {
	Host 			string
	Port			string
	User 			string
	Password 		string
	Dbname 			string
	Sslmode 		string
}

type Redis struct {
	Host			string
	Port			string
	Password		string
	Db				int
}

type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

func LoadConfig (filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("Config file now found")
		}
	}	

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}