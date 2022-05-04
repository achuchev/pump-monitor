package main

import (
	"os"

	"github.com/achuchev/pump-monitor/cmd/command"
	"github.com/achuchev/pump-monitor/cmd/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer panicHandler()
	loggingLevel := log.InfoLevel
	envDebug := os.Getenv("PUMP_DEBUG")

	switch envDebug {
	case "1":
		{
			loggingLevel = log.DebugLevel
		}
	case "2":
		{
			loggingLevel = log.TraceLevel
		}
	}
	logger.LoggerInit(loggingLevel)

	if err := command.RootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}

func panicHandler() int {
	var err error
	if r := recover(); r != nil {
		log.Errorf("panic : %v", err)
	}
	return 0
}
