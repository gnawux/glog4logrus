package glog4logrus

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/golang/glog"
)

type GlogFormatter struct {}

func (f *GlogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.appendValue(b, false, fmt.Sprintf("%c%-44s", strings.ToUpper(entry.Level.String())[0], entry.Message))
	for _, key := range keys {
		f.appendKeyValue(b, key, entry.Data[key])
	}

	return b.Bytes(), nil
}

func (f *GlogFormatter) needsQuoting(text string) bool {
	if len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *GlogFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, true, value)
}

func (f *GlogFormatter) appendValue(b *bytes.Buffer, quote bool, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !quote || !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

type GlogOuptut struct {}

// Through the call stack, why we define baseDepth as 3 here:
//   1: the Write() here
//   2: entry.go: Entry.log()
// * 3: entry.go: Entry.{Debug, Info, ...}
//   4: entry.go: Entry.{Debugf, Infof, ...}; Entry.{Debugln, Infoln, ...}; Entry.{Print, Warning}; logger.go/exporter.go Logger.{Debug, Info,...}; writer.go: writerScanner
//   5: entry.go: Entry.{Println, Printf, Warningf, Warningln}; logger.go/exporter.go: Logger.{Debugf...}; Logger.{Debugln...}; Logger.{Print, Warning}
//   6: logger.go/exporter.go: Logger.{Println, Printf, Warningf, Warningln}
// And if we met `writerScanner()`, it means we cannot find the original post because it is in an independent goroutine from the original caller, we should giv up.
const baseDepth = 3
var logFunctions = map[string]bool {
	"Debug": true,
	"Info": true,
	"Print": true,
	"Warn": true,
	"Warning": true,
	"Error": true,
	"Fatal": true,
	"Panic": true,
	"Debugf": true,
	"Infof": true,
	"Printf": true,
	"Warnf": true,
	"Warningf": true,
	"Errorf": true,
	"Fatalf": true,
	"Panicf": true,
	"Debugln": true,
	"Infoln": true,
	"Println": true,
	"Warnln": true,
	"Warningln": true,
	"Errorln": true,
	"Fatalln": true,
	"Panicln": true,
}

func (w *GlogOuptut) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	pc := make([]uintptr, 5)
	depth := baseDepth
	stacks := runtime.Callers(depth+1, pc)
	if stacks != 0 {
		frames := runtime.CallersFrames(pc)
		next := true
		for next {
			frame, more := frames.Next()
			if strings.HasPrefix(frame.Function, "github.com/sirupsen/logrus.") &&
				frame.Function != "github.com/sirupsen/logrus.(*Entry).writerScanner" {
				//still in logger, skip this caller
				depth++
				next = more
			} else {
				next = false
			}
		}
	}

	msg := string(p)
	switch msg[0] {
	case 'D', 'I':
		glog.InfoDepth(depth, msg[1:])
	case 'W':
		glog.WarningDepth(depth, msg[1:])
	case 'E':
		glog.ErrorDepth(depth, msg[1:])
	case 'F':
		glog.FatalDepth(depth, msg[1:])
	case 'P':
		glog.FatalDepth(depth, msg[1:])
	default:
		glog.InfoDepth(depth, msg[1:])
	}
	return len(p),nil
}

func GlogLevel() logrus.Level {
	if glog.V(1) {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}
