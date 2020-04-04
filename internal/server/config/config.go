package config

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/util"

	"github.com/spf13/viper"
)

var (
	Cfg map[string]string
	log *logrus.Logger
)

func init() {
	log = util.GetLogger()
}

// LoadConfigFile ... Load the config file for the server
func LoadConfigFile(cfgFile string) {
	log.Debugln("Input to LoadConfigFile function is :: ", cfgFile)
	serverConfig := loadConfig(cfgFile)

	if len(serverConfig) > 0 {
		Cfg = serverConfig
	}
}

func readConfig(filename, dirname string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(dirname)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

// Load the config as default or from config file
func loadConfig(cfgFileName string) map[string]string {
	serverConfig := make(map[string]string)

	filename := strings.TrimSuffix(filepath.Base(cfgFileName), filepath.Ext(filepath.Base(cfgFileName)))
	dirname := filepath.Dir(cfgFileName)

	log.Debugln("Directory of config file is :: ", dirname)
	log.Debugln("Filename of config file is :: ", filename)

	v1, err := readConfig(filename, dirname, map[string]interface{}{
		"server": map[string]interface{}{
			"address": "0.0.0.0",
			"port":    9090,
		},
		"http": map[string]interface{}{
			"read.timeout":   0,
			"write.timeout":  0,
			"maxrequestsize": 999999,
		},
		"logging": map[string]interface{}{
			"level": "debug",
		},
		"apps": map[string]interface{}{
			"directory": "apps",
		},
	})

	if err != nil {
		panic(fmt.Errorf("Error when reading config %v ", err))
	}
	log.Debugln("Configuration is :: ", v1)

	serverConfig["server.port"] = strconv.Itoa(v1.Get("server.port").(int))
	serverConfig["server.address"] = v1.Get("server.address").(string)
	serverConfig["logging.level"] = v1.Get("logging.level").(string)
	serverConfig["timeout.read"] = strconv.Itoa(v1.Get("http.read.timeout").(int))
	serverConfig["timeout.write"] = strconv.Itoa(v1.Get("http.write.timeout").(int))
	serverConfig["http.maxrequestsize"] = strconv.Itoa(v1.Get("http.maxrequestsize").(int))
	serverConfig["apps.directory"] = v1.Get("apps.directory").(string)

	log.Debugln("The config map created is :: ", serverConfig)
	return serverConfig
}
