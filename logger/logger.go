package logger

import (
	"net"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	conn, err := net.Dial("tcp", "logstash:5000")
	if err != nil {
		Log.Fatalf("Failed to connect to Logstash: %v", err)
	}
	Log.Out = conn
	Log.Formatter = &logrus.JSONFormatter{}
}
