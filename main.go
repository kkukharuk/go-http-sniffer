package main

import (
	"flag"
	"fmt"
	"git.kukharuk.ru/kkukharuk/go-http-sniffer/logger"
	"git.kukharuk.ru/kkukharuk/go-http-sniffer/router"
	"net/http"
	"os"
)

var (
	host        string
	port        int
	enableTls   bool
	certificate string
	key         string
	targetUrl   string
	logFile     string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "Listen sniffer host")
	flag.IntVar(&port, "port", 8080, "Listen sniffer port")
	flag.BoolVar(&enableTls, "tls", false, "Enable TLS")
	flag.StringVar(&certificate, "certfile", "./file.cer", "Certificate file")
	flag.StringVar(&key, "keyfile", "./file.key", "Key file")
	flag.StringVar(&targetUrl, "target", "http://127.0.0.1:5000", "Target url")
	flag.StringVar(&logFile, "logFile", "./httpSniffer-requests.log", "Logfile")
	flag.Parse()
	_, errCert := os.Stat(certificate)
	_, errKey := os.Stat(key)
	if enableTls {
		if errCert != nil && errKey == nil {
			logger.Fatal(fmt.Sprintf("'%s': File not found", certificate))
		} else if errCert == nil && errKey != nil {
			logger.Fatal(fmt.Sprintf("'%s': File not found", key))
		} else if errCert != nil && errKey != nil {
			logger.Fatal(fmt.Sprintf("'%s' and ''%s: File not found", certificate, key))
		}
		if port == 8080 {
			port = 8443
		}
	}
	logger.Init("Go-HTTP-Sniffer", "0.0.1-betta", logger.LogDetail, logger.LogLevelDebug)
}

func main() {
	r := router.New(targetUrl, logFile)
	if enableTls {
		logger.Info(fmt.Sprintf("Running server on htts://%s:%d", host, port))
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", host, port), certificate, key, r.RootHandler()); err != nil {
			panic(err)
		}
	} else {
		logger.Info(fmt.Sprintf("Running server on http://%s:%d", host, port))
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r.RootHandler()); err != nil {
			panic(err)
		}
	}
}
