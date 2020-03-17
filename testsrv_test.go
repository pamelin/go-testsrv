package testsrv

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	srv, rec := NewServer(h)
	require.NotNil(t, srv)
	require.NotNil(t, rec)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPut, srv.URL+"/test-path", bytes.NewBuffer([]byte("hello")))
	require.NoError(t, err)
	c := http.Client{Timeout: 2 * time.Second}

	resp, err := c.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, http.MethodPut, rec.Req.Method)
	assert.Equal(t, "/test-path", rec.Req.URL.String())
	assert.Equal(t, "hello", string(rec.Body))
}

func TestNewServerWithStatus(t *testing.T) {
	srv, rec := NewServerWithStatus(http.StatusNotFound)
	require.NotNil(t, srv)
	require.NotNil(t, rec)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/my/test-path", nil)
	require.NoError(t, err)
	c := http.Client{Timeout: 2 * time.Second}

	resp, err := c.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotNil(t, rec.Req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, http.MethodGet, rec.Req.Method)
	assert.Equal(t, "/my/test-path", rec.Req.URL.String())
	assert.Empty(t, rec.Body)
}

func TestNewServerWithBody(t *testing.T) {
	srv, rec := NewServerWithBody("yo!!!")
	require.NotNil(t, srv)
	require.NotNil(t, rec)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/kindof-test-path", nil)
	require.NoError(t, err)
	c := http.Client{Timeout: 2 * time.Second}

	resp, err := c.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "yo!!!", string(body))

	assert.Equal(t, http.MethodGet, rec.Req.Method)
	assert.Equal(t, "/kindof-test-path", rec.Req.URL.String())
	assert.Empty(t, rec.Body)
}
