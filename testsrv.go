package testsrv

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

func NewServer(f http.HandlerFunc) (*httptest.Server, *ReqRec) {
	rec := &ReqRec{}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		rec.record(r)
		if f != nil {
			f(w, r)
		}
	})
	return httptest.NewServer(h), rec
}

func NewServerWithStatus(status int) (*httptest.Server, *ReqRec) {
	return NewServer(NewHandler(status, nil))
}

func NewServerWithBody(body string) (*httptest.Server, *ReqRec) {
	return NewServer(NewHandler(http.StatusOK, []byte(body)))
}

type ReqRec struct {
	Req  *http.Request
	Body []byte
}

func (rec *ReqRec) record(r *http.Request) {
	var body io.Reader = nil
	if b, err := ioutil.ReadAll(r.Body); err != nil {
		log.Fatalf("error while reading body: %v", err)
	} else {
		rec.Body = b
	}
	rec.Req = httptest.NewRequest(r.Method, r.URL.String(), body)
	rec.Req.Header = r.Header
}

func NewHandler(status int, body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			if _, err := w.Write(body); err != nil {
				log.Fatalf("Error writing response body: %v", err)
			}
		}
	}
}
