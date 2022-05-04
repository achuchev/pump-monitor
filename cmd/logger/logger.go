package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"golang.org/x/term"
)

func LoggerInit(level log.Level, verbose ...bool) {
	log.StandardLogger().ReplaceHooks(make(log.LevelHooks))
	log.SetLevel(level)

	logDirPath := "./"
	if runtime.GOOS == "linux" {
		logDirPath = "/var/log/"
	}

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logDirPath + "heat_pump_info.log",
		logrus.ErrorLevel: logDirPath + "heat_pump_error.log",
	}
	log.AddHook(lfshook.NewHook(pathMap, &logrus.TextFormatter{}))

	if level == log.DebugLevel || level == log.TraceLevel {
		var debugFormatter = &zt_formatter.ZtFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
		log.SetReportCaller(true)
		log.SetFormatter(debugFormatter)
	} else {
		forceColors := false
		if isTerminal(log.StandardLogger().Out) && runtime.GOOS == "windows" {
			forceColors = true
			log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
		}
		//DisableColors:          !isTerminal(os.Stderr),
		log.SetFormatter(&log.TextFormatter{
			ForceColors:            forceColors,
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	}
	log.Trace("Logger initialized")
	log.Tracef("Terminal Detected: %v", isTerminal(os.Stderr))
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
