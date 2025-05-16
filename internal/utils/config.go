package utils

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBSource           string        `mapstructure:"DB_SOURCE"`      //db !!
	ServerPort         string        `mapstructure:"SERVER_PORT"`    //server !!
	ServerShutDownTime time.Duration `mapstructure:"SHUTDOWN_TIME"`  //server !!
	MOCK1CAddress      string        `mapstructure:"MOCK1C_ADDRESS"` //client !!
	CronBatchSize      int           `mapstructure:"BATCH_SIZE"`     //cron !!
	CronWorkerCount    int           `mapstructure:"WORKER_COUNT"`   //cron !!
	WorkerTimeout      time.Duration `mapstructure:"WORKER_TIMEOUT"` //eno
	MaxRetries         int32         `mapstructure:"MAX_RETRIES"`    //eno
	JobInterval        time.Duration `mapstructure:"JOB_INTERVAL"`   //eno
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)

	return
}
