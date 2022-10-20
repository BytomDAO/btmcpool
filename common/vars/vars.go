package vars

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to yml config file")
	flag.Parse()

	if len(configPath) == 0 {
		return
	}

	initImpl(configPath)
}

func initImpl(configPath string) {
	ext := path.Ext(configPath)
	name := strings.Replace(path.Base(configPath), ext, "", -1)
	dir := path.Dir(configPath)

	// name of config file (without extension)
	viper.SetConfigName(name)
	// look for config in the working directory
	viper.AddConfigPath(dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file changed: %v", e.Name)
	})
}

func GetBool(name string, def bool) bool {
	viper.SetDefault(name, def)
	return viper.GetBool(name)
}

func GetInt64(name string, def float64) int64 {
	viper.SetDefault(name, def)
	return viper.GetInt64(name)
}

func GetInt(name string, def float64) int {
	viper.SetDefault(name, def)
	return viper.GetInt(name)
}

func GetFloat64(name string, def float64) float64 {
	viper.SetDefault(name, def)
	return viper.GetFloat64(name)
}

func GetString(name string, def string) string {
	viper.SetDefault(name, def)
	return viper.GetString(name)
}

func GetTime(name string, def time.Time) time.Time {
	viper.SetDefault(name, def)
	return viper.GetTime(name)
}

func GetDuration(name string, def time.Duration) time.Duration {
	viper.SetDefault(name, def)
	return viper.GetDuration(name)
}

func GetStringSlice(name string, def []string) []string {
	viper.SetDefault(name, def)
	return viper.GetStringSlice(name)
}

// Unmarshal unmarshal config value to config struct
func Unmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}

// Fill fills a struct with config values under given sub-tree
func Fill(name string, cfg interface{}) error {
	sub := viper.Sub(name)
	if sub == nil {
		return errors.New("No Such name: " + name)
	}

	if err := sub.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}
