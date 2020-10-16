package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init 初始化
func Init() {
	if lvl, err := logrus.ParseLevel(viper.GetString("log.level")); err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(lvl)
	}
	logrus.SetReportCaller(viper.GetBool("log.report_caller"))
}
