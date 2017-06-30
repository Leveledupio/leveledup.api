package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	app "github.com/strongjz/leveledup-api/application"
	"gopkg.in/op/go-logging.v1"
	"os"
)

var (
	log = logging.MustGetLogger("gbs")

	defaultFormat = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
	//ENV sets the environment so the right config is loaded
	ENV = "dev"
)

//Init function sets the default env to dev
func init() {
	log.Debugf("Init: Default ENV: %v", ENV)

	if os.Getenv("ENV") != "" {
		ENV = os.Getenv("ENV")
	}

	log.Debugf("Init: ENV Set: %v", ENV)

}

//newConfig function parsers the config file for database connection etc
//
func newConfig() (*viper.Viper, error) {

	configName := fmt.Sprintf("config/%v-config.yaml", ENV)

	log.Debugf("Loading Config file %v", configName)

	c := viper.New()
	c.AddConfigPath(".")
	c.SetConfigFile(configName)
	c.SetConfigType("yaml")
	//c.Debug()
	c.WatchConfig()

	c.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed:", e.Name)
	})

	c.AutomaticEnv()

	err := c.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	c.Debug()

	return c, nil
}

func main() {

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend,
			logging.MustStringFormatter(defaultFormat),
		),
	)

	config, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Config: %v", config)
	log.Debug("App DSN: %s", a.DSN)

	r := RouteSetup(a)

	r.Run() // listen and serve on 0.0.0.0:8080
}
