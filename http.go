package http

import (
	"bytes"
	"io"
	"net/http"

	"storj.io/drpc"
)

var _ drpc.Transport = &httpTransport{}

type httpTransport struct {
	c        *http.Client
	method   string
	endpoint string
	ct       string
	closed   bool
	body     io.ReadCloser
}

func NewTransport(c *http.Client, endpoint, method, ct string) *httpTransport {
	if ct == "" {
		ct = "application/drpc"
	}
	return &httpTransport{c: c, method: method, endpoint: endpoint, ct: ct}
}

func (t *httpTransport) Read(buf []byte) (int, error) {
	return t.body.Read(buf)
}

func (t *httpTransport) Write(buf []byte) (n int, err error) {
	var req *http.Request
	var rsp *http.Response

	req, err = http.NewRequest(t.method, t.endpoint, bytes.NewReader(buf))
	if err != nil {
		return n, err
	}
	req.Header.Add("Content-Type", t.ct)
	rsp, err = t.c.Do(req)
	if err != nil {
		return n, err
	}

	t.body = rsp.Body
	return len(buf), err
}

func (t *httpTransport) Close() error {
	if t.body == nil {
		return nil
	}
	return t.body.Close()
}
