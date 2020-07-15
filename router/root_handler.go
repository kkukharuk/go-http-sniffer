package router

import (
	"fmt"
	"git.kukharuk.ru/kkukharuk/go-http-sniffer/logger"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type rootHandler struct {
	targetUrl string
	logFile   string
	client    *http.Client
}

func (h rootHandler) httpClient(r *http.Request, requestBody []byte) (*http.Response, []byte, error) {
	var req *http.Request
	var err error
	var responseBody []byte
	var newFileData string
	uuid := r.Header.Get("Request-ID")
	r.Header.Del("Request-ID")
	resTmpl := `<= Sesponse (UUID: %s)
   Response:
     ResponseCode: %d
     Headers:
%s
     Body: %s

`
	oldFileData, err := ioutil.ReadFile(h.logFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Error read log-file: %s", err.Error()))
	}
	newFileData += string(oldFileData)
	if len(requestBody) != 0 {
		req, err = http.NewRequest(r.Method, h.targetUrl+r.RequestURI, strings.NewReader(string(requestBody)))
	} else {
		req, err = http.NewRequest(r.Method, h.targetUrl+r.RequestURI, nil)
	}
	if err != nil {
		return nil, responseBody, err
	}
	for headerName, _ := range r.Header {
		if headerName == "Accept" {
			req.Header.Add("Accept", "applications/json")
		} else {
			req.Header.Add(headerName, r.Header.Get(headerName))
		}
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, responseBody, err
	}
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, responseBody, err
	}
	defer resp.Body.Close()
	_ = ioutil.WriteFile("test/test.json", responseBody, 0766)
	var headers string
	headersCount := len(resp.Header)
	i := 1
	for headerKey, headerValue := range resp.Header {
		headers += fmt.Sprintf("       %q: %q", headerKey, headerValue)
		if i != headersCount {
			headers += "\n"
		}
		i++
	}
	newFileData += fmt.Sprintf(resTmpl, uuid, resp.StatusCode, headers, string(responseBody))
	err = ioutil.WriteFile(h.logFile, []byte(newFileData), 0766)
	if err != nil {
		logger.Error(fmt.Sprintf("Error write request info to log-file: %s", err.Error()))
	} else {
		logger.Debug("Write request info to log-file is successfully complite")
	}
	return resp, responseBody, nil
}

func (h rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New()
	var newFileData string
	resTmpl := `=> Request (UUID: %s)
   Request:
     Request URI: %s
     Method: %s
     Headers:
%s
     Body: %s
`
	if _, err := os.Stat(h.logFile); os.IsNotExist(err) {
		_, err = os.Create(h.logFile)
	}
	oldFileData, err := ioutil.ReadFile(h.logFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Error read log-file: %s", err.Error()))
	}
	newFileData += string(oldFileData)
	var headers string
	headersCount := len(r.Header)
	i := 1
	for headerKey, headerValue := range r.Header {
		headers += fmt.Sprintf("       %q: %q", headerKey, headerValue)
		if i != headersCount {
			headers += "\n"
		}
		i++
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("Error read request body: %s", err.Error()))
	}
	newFileData += fmt.Sprintf(resTmpl, uuid.String(), r.RequestURI, r.Method, headers, string(requestBody))
	err = ioutil.WriteFile(h.logFile, []byte(newFileData), 0766)
	if err != nil {
		logger.Error(fmt.Sprintf("Error write request info to log-file: %s", err.Error()))
	} else {
		logger.Debug("Write request info to log-file is successfully complite")
	}
	r.Header.Add("Request-ID", uuid.String())
	var targetResponse *http.Response
	var responseBody []byte
	if len(requestBody) != 0 {
		targetResponse, responseBody, err = h.httpClient(r, requestBody)
	} else {
		targetResponse, responseBody, err = h.httpClient(r, []byte{})
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Error target request: %s", err.Error()))
	}
	for headerName, _ := range targetResponse.Header {
		if headerName == "Location" {
			var replceUrl string
			if r.TLS == nil {
				replceUrl = fmt.Sprintf("http://%s", r.Host)
			} else {
				replceUrl = fmt.Sprintf("https://%s", r.Host)
			}
			w.Header().Set(headerName, strings.Replace(targetResponse.Header.Get(headerName),
				h.targetUrl,
				replceUrl,
				-1))
		} else {
			w.Header().Set(headerName, targetResponse.Header.Get(headerName))
		}
	}
	if _, err = io.Copy(w, strings.NewReader(string(responseBody))); err != nil {
		logger.Error(fmt.Sprintf("Error write response body: %s", err.Error()))
	}
	w.WriteHeader(targetResponse.StatusCode)
}
