package config

import "github.com/spf13/viper"

type Config struct {
	KafkaURL        string `mapstructure:"KAFKA_URL"`
	MongoURL        string `mapstructure:"MONGO_URL"`
	PostgresURL     string `mapstructure:"POSTGRES_URL"`
	APIGWport       string `mapstructure:"APIGW_PORT"`
	ProtocolVersion string `mapstructure:"PROTOCOL_VERSION"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
