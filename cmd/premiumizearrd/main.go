package main

import (
	"flag"

	"github.com/jackdallas/premiumizearr/internal/utils"
)

func main() {
	//Flags
	var logLevel string
	var configFile string
	var loggingDirectory string

	//Parse flags
	flag.StringVar(&logLevel, "log", utils.EnvOrDefault("PREMIUMIZEARR_LOG_LEVEL", "info"), "Logging level: \n \tinfo,debug,trace")
	flag.StringVar(&configFile, "config", utils.EnvOrDefault("PREMIUMIZEARR_CONFIG_DIR_PATH", "./"), "The directory the config.yml is located in")
	flag.StringVar(&loggingDirectory, "logging-dir", utils.EnvOrDefault("PREMIUMIZEARR_LOGGING_DIR_PATH", "./"), "The directory logs are to be written to")
	flag.Parse()

	App := &App{}
	App.Start(logLevel, configFile, loggingDirectory)

}
