package g

import (
	log "github.com/Sirupsen/logrus"
)

func InitLog(level string) (err error) {
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.Fatal("log conf only allow [info, debug, warn, error, fatal, panic], please check your confguire")
	}
	return
}

func IsDebug() bool {
	return log.GetLevel() == log.DebugLevel
}
