package glog4logrus

import (
	"flag"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/golang/glog"
)

func TestBasic(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Parse()

	logrus.SetLevel(GlogLevel())
	logrus.SetFormatter(&GlogFormatter{})
	logrus.SetOutput(&GlogOuptut{})

	glog.Info("this is from glog")

	logrus.Debug("this is debug")
	logrus.Info("this is info")
	logrus.Printf("this is a %s", "printf")

	xlogger := logrus.WithField("logger", "xlogger")
	xlogger.Info("this is info")
	xlogger.WithField("special", "should be quote").Info("this is info again")

	xlogger.WithField("annotation", "long line").Info("here is a long line, it should be very very very very long, longer than the magic 44 characters")

	logrus.Info("")
	glog.Flush()
}
