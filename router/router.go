package router

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Router struct {
	targetUrl   string
	logFile     string
	client      *http.Client
	rootHandler rootHandler
}

func New(targetUrl, file string) *Router {
	return &Router{
		targetUrl: targetUrl,
		logFile:   file,
		client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (r *Router) RootHandler() http.Handler {
	h := r.rootHandler
	h.targetUrl = r.targetUrl
	h.logFile = r.logFile
	h.client = r.client
	return h
}
