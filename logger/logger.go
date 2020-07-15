package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	basePath      = "git.kukharuk.ru/kkukharuk/go-http-sniffer"
	LogDetail     = "DETAIL"
	LogShort      = "SHORT"
	LogVersion    = "0.1.0"
	InitMsgLength = 47
	LogLevelInfo  = "INFO"
	LogLevelDebug = "DEBUG"
	LogLevelError = "ERROR"
	LogLevelWarn  = "WARN"
)

type logFlags struct {
	Level string
}

var (
	basePathPackagesDel = ""
	loggerFlags         = logFlags{Level: LogLevelError}
)

type LogWriterWithTimestamp struct {
}

func (writer LogWriterWithTimestamp) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + " " + strings.Replace(string(bytes), "???", "", -1))
}

type LogWriterWithoutTimestamp struct {
}

func (writer LogWriterWithoutTimestamp) Write(bytes []byte) (int, error) {
	return fmt.Print(strings.Replace(string(bytes), "???", "", -1))
}

type LogWriterDetailed struct {
}

func (writer LogWriterDetailed) Write(bytes []byte) (int, error) {
	_, gfn, line, _ := runtime.Caller(4)
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + " " + strings.Replace(string(bytes), "???", fmt.Sprintf(" %s:%d ::", strings.Replace(strings.Replace(strings.Replace(gfn, basePathPackagesDel, "", -1), "/", ".", -1), ".go", "", -1), line), -1))
	//return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + " " + fmt.Sprintf(string(bytes), strings.Replace(strings.Replace(strings.Replace(gfn, basePathPackagesDel, "", -1), "/", ".", -1), ".go", "", -1), line))
}

func Init(appId string, appVersion string, logType string, logLevel string) {
	log.SetFlags(0)

	setLogLevel(logLevel)

	basePathPackagesDel = fmt.Sprintf("%s/src/%s/", strings.Replace(os.Getenv("GOPATH"), "\\", "/", -1), basePath)

	if logType == LogDetail {
		log.SetOutput(new(LogWriterDetailed))
	} else if logType == LogShort {
		log.SetOutput(new(LogWriterWithoutTimestamp))
	} else {
		log.SetOutput(new(LogWriterWithTimestamp))
	}

	log.Println("***********************************************")
	log.Println("***********************************************")
	log.Println("***                                         ***")
	log.Println("***          (''').o___o.(''')   П          ***")
	log.Println("***           \\  '' o_o ''  /    Р          ***")
	log.Println("***            \\   \\_Ш_/  /      Е          ***")
	log.Println("***             |         |      В          ***")
	log.Println("***             /   /U\\   \\      Е          ***")
	log.Println("***            (,,,)   (,,,)     Д          ***")
	log.Println("***                                         ***")
	log.Println("***********************************************")
	log.Println("***********************************************")
	msg := fmt.Sprintf("***  Application: %s  ***", appId)
	if len(msg) < InitMsgLength {
		msg = fmt.Sprintf("***  Application: %s", appId)
		r := InitMsgLength - len(msg)
		for r != 5 {
			msg += " "
			r -= 1
		}
		msg += "  ***"
	}
	log.Println(msg)
	msg = fmt.Sprintf("***  Version: %s  ***", appVersion)
	if len(msg) < InitMsgLength {
		msg = fmt.Sprintf("***  Version: %s", appVersion)
		r := InitMsgLength - len(msg)
		for r != 5 {
			msg += " "
			r -= 1
		}
		msg += "  ***"
	}
	log.Println(msg)
	log.Println(fmt.Sprintf("***  LoggerVersion: %s                   ***", LogVersion))
	log.Println("***********************************************")
	log.Println("***********************************************")
}

func Info(message string) {
	if loggerFlags.Level == LogLevelInfo || loggerFlags.Level == LogLevelDebug {
		log.Println("[ INFO ]??? " + message)
	}
}

func Warn(message string) {
	if loggerFlags.Level == LogLevelInfo || loggerFlags.Level == LogLevelWarn || loggerFlags.Level == LogLevelDebug {
		log.Println("[ WARN ]??? " + message)
	}
}

func Error(message string) {
	if loggerFlags.Level == LogLevelError || loggerFlags.Level == LogLevelWarn || loggerFlags.Level == LogLevelError || loggerFlags.Level == LogLevelDebug {
		log.Println("[ ERROR ]??? " + message)
	}
}

func Debug(message string) {
	if loggerFlags.Level == LogLevelDebug {
		log.Println("[ DEBUG ]??? " + message)
	}
}

func Fatal(message string) {
	log.Fatalln("[ FATAL ]??? " + message)
}

func setLogLevel(level string) {
	if level != "" && (level == LogLevelError || level == LogLevelWarn || level == LogLevelInfo || level == LogLevelDebug) {
		loggerFlags.Level = level
	} else if level != "" {
		Fatal("Error parse logger level flag")
	}
}
