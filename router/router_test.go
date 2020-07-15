package router_test

import (
	"git.kukharuk.ru/kkukharuk/go-http-sniffer/router"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	r := router.New()
	srv := httptest.NewServer(r.RootHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/opa")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
