package logger

import "github.com/sirupsen/logrus"

func New(level, format string) *logrus.Logger {
	log := logrus.New()
	parsed, err := logrus.ParseLevel(level)
	if err == nil {
		log.SetLevel(parsed)
	}
	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
	log.SetReportCaller(true)
	return log
}
