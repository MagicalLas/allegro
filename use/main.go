package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gomod.usaken.org/use/config"
)

var c = &config.Config{}

func init() {
	currentProcessLocation, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("fail to load current bin location for load default config file location: %w", err))
	}

	defaultConfigFilePath := filepath.Dir(currentProcessLocation) + "/.use.config"
	configFilePath := flag.String("configpath", defaultConfigFilePath, "config file path for other config file location or load other config file.")
	newConfig := flag.Bool("newconfig", true, "if don't exist config file in location, create default empty config file for further setting.")

	c, err = config.LoadFromFile(*configFilePath, *newConfig)
	if err != nil {
		panic(fmt.Errorf("fail to load config file: %w", err))
	}
}

func main() {

}
