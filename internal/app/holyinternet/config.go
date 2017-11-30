package holyinternet

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	appName = "holy-internet"
	checkIntervalInSeconds = 3
	checkHostsCount = 3
)

var configPaths = []string{
	fmt.Sprintf("/etc/%s/", appName),
	fmt.Sprintf("$HOME/.%s", appName),
	".",
}

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetDefault("pray.every", checkIntervalInSeconds)
	viper.SetDefault("pray.count", checkHostsCount)

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	angels := viper.GetStringSlice("saints")
	if len(angels) == 0 {
		panic("No saints found in config file! Have a little faith!")
	}
}
