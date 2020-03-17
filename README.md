# go-testsrv
Helper for the golang's httptest which spins up simple server with a http response or body and 
provides a request recorder which can be used in your assertions.

## Example usage

```
func TestNewServerWithStatus(t *testing.T) {
	srv, rec := NewServerWithStatus(http.StatusNotFound)
	defer srv.Close()

	resp, err := c.Get(srv.URL + "/my/test-path")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    require.NotNil(t, rec.Req)
	assert.Equal(t, http.MethodGet, rec.Req.Method)
	assert.Equal(t, "/my/test-path", rec.Req.URL.String())
	assert.Empty(t, rec.Body)
}
```
