package configuration

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configFilePath = "/var/goaway/config.yml"
	goAwayConfig   *viper.Viper

	once sync.Once
)

// GetInstance returns an instance of viper config
func GetInstance() *viper.Viper {
	once.Do(func() {
		loadConfig()
	})
	return goAwayConfig
}

func loadConfig() {
	config := viper.New()

	// Setting defaults for this application
	config.SetDefault("addr", ":10000")
	config.SetDefault("debug", true)

	config.BindEnv("node.ip", "HOST_IP")

	config.SetConfigFile(configFilePath)

	config.OnConfigChange(OnConfigChanged)
	config.WatchConfig()

	err := config.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("can't read config file")
		goAwayConfig = config
		return
	}
	logrus.Infof("configuration file is loaded from %s", configFilePath)
	SetDebugLogLevel(config.GetBool("debug"))

	logrus.Debugf("loaded config: %v", config.AllSettings())
	goAwayConfig = config
}

// SetFilePath sets path of config file
func SetFilePath(filePath string) {
	if filePath != configFilePath {
		configFilePath = filePath
		loadConfig()
	}
}

// SetDebugLogLevel sets log level to debug mode
func SetDebugLogLevel(isDebug bool) {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("log level is set to Debug")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// OnConfigChanged excuates when config changes
func OnConfigChanged(_ fsnotify.Event) {
	loadConfig()
	logrus.Info("configuration is reloaded")
}
