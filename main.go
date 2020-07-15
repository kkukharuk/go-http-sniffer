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
	enableTls   bool
	certificate string
	key         string
	targetUrl   string
	logFile     string
)

func init() {
	logger.Init("Go-HTTP-Sniffer", "0.0.1-betta", logger.LogDetail, logger.LogLevelDebug)
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
	}
}

func main() {
	r := router.New(targetUrl, logFile)
	if enableTls {
		if err := http.ListenAndServeTLS(":8443", certificate, key, r.RootHandler()); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(":8080", r.RootHandler()); err != nil {
			panic(err)
		}
	}
}
