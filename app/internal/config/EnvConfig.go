package config

import (
	"log"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

func init() {
	if _, err := os.Stat(".env"); err == nil {
		EnvConfigs = loadFromEnvFile()
	} else {
		EnvConfigs = loadEnvVariables()
	}
}

type envConfigs struct {
	Port              int    `mapstructure:"PORT"`
	CognitoClientId   string `mapstructure:"COGNITO_CLIENT_ID"`
	CognitoUserPoolId string `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoRegion     string `mapstructure:"COGNITO_REGION"`
	AppEnv            string `mapstructure:"APP_ENV"`
}

func loadEnvVariables() *envConfigs {
	viper.AutomaticEnv()

	config := &envConfigs{
		Port:              viper.GetInt("PORT"),
		CognitoClientId:   viper.GetString("COGNITO_CLIENT_ID"),
		CognitoUserPoolId: viper.GetString("COGNITO_USER_POOL_ID"),
		CognitoRegion:     viper.GetString("COGNITO_REGION"),
		AppEnv:            viper.GetString("APP_ENV"),
	}

	validateEnvVariables(config)
	return config
}

func loadFromEnvFile() *envConfigs {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	config := &envConfigs{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal(err)
	}

	validateEnvVariables(config)
	return config
}

func validateEnvVariables(config *envConfigs) {
	configType := reflect.TypeOf(*config)
	configValue := reflect.ValueOf(*config)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envName := field.Tag.Get("mapstructure")
		expectedType := field.Type.Name()

		switch expectedType {
		case "int":
			if value := configValue.Field(i).Int(); value == 0 {
				log.Fatalf("Environment variable %s is not set", envName)
			}
		case "string":
			if value := configValue.Field(i).String(); value == "" {
				log.Fatalf("Environment variable %s is not set", envName)
			}
		//TODO: Add bool type
		// case "bool":
		default:
			log.Fatalf("Environment variable %s is not set", envName)

		}
	}
}
