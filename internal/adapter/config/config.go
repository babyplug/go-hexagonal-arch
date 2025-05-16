package config

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var (
	_config Config = Config{
		MongoURI: "mongodb://localhost:27017",

		Port:     ":8080",
		Duration: "86400s", // 24 hours
	}
)

type Config struct {
	MongoURI       string `mapstructure:"MONGO_URI"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	Port           string `mapstructure:"PORT"`
	Env            string `mapstructure:"ENV"`
	AllowedOrigins string `mapstructure:"ALLOWED_ORIGINS"`
	AllowedMethods string `mapstructure:"ALLOWED_METHODS"`
	AllowedHeaders string `mapstructure:"ALLOWED_HEADERS"`
	Duration       string `mapstructure:"DURATION"`
}

func Load() *Config {
	viper := viper.New()

	configFile, ok := os.LookupEnv("CONFIG_FILE")
	if ok && len(configFile) > 0 {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Failed to read config file: " + err.Error())
			viper.AutomaticEnv()
		}
	} else {
		viper.AutomaticEnv()
	}

	bindEnv(_config)

	err := viper.Unmarshal(&_config)
	if err != nil {
		log.Printf("Failed to unmarshal config: " + err.Error())
	}

	log.Printf("Config: %+v\n", _config)

	return &_config
}

func bindEnv(dest any, parts ...string) {
	ifv := reflect.ValueOf(dest)
	ift := reflect.TypeOf(dest)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)

		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			bindEnv(v.Interface(), append(parts, tv)...)
		default:
			envKey := strings.Join(append(parts, tv), ".")
			err := viper.BindEnv(envKey)
			if err != nil {
				log.Printf("bind env key %s failed: %v\n", envKey, err)
			}
		}
	}
}
