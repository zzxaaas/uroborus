package settings

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"path"
	"runtime"
	"uroborus/common"
)

const configDir = "config"

// Config -
type Config struct {
	Postgres struct {
		Host      string `mapstructure:"host"`
		DBName    string `mapstructure:"dbname"`
		Port      int    `mapstructure:"port"`
		Username  string `mapstructure:"user"`
		Password  string `mapstructure:"passwd"`
		BatchSize int    `mapstructure:"batchSize"`
	} `mapstructure:"postgres"`

	Env       string `mapstructure:"env"`
	GinMode   string `mapstructure:"gin_mode"`
	ApiPrefix string `mapstructure:"apiPrefix"`
}

var (
	currentStaticCompiledAbsFilename string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	currentStaticCompiledAbsFilename = filename
}

// GetConfig -
func GetConfig() *Config {
	viper.SetConfigName("conf")
	viper.AddConfigPath("/" + configDir)
	viper.AddConfigPath(configDir) // first load global config
	// !!!Important, this line config path is compile-time path
	viper.AddConfigPath(path.Join(currentStaticCompiledAbsFilename, "..", "..", configDir))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatal("fail to load config file", err)
	}
	{ // set defaults
		viper.SetDefault("taskCheckTTL", "15s")
		viper.SetDefault("gin_mode", gin.DebugMode)
		viper.SetDefault("postgresConfig.batchSize", 1000)
		viper.SetDefault("maxQueryWindowNum", 100)
	}
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("unmarshal config to struct failed!", err)
	}
	if viper.GetBool("debug") {
		if configJson, err := json.MarshalIndent(config, "", "  "); err != nil {
			panic(err)
		} else {
			log.Printf("%+v\n", string(configJson))
		}
	}
	return config
}

// InDebugMode -
func (config *Config) InDebugMode() bool {
	return config.GinMode != common.ReleaseMode && config.GinMode != common.TestMode
}

// InReleaseMode -
func (config *Config) InReleaseMode() bool {
	return config.GinMode == common.ReleaseMode
}

// InTestMode -
func (config *Config) InTestMode() bool {
	return config.GinMode == common.TestMode
}
