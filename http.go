package http

import (
	"bytes"
	"io/ioutil"
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
}

func NewTransport(c *http.Client, endpoint, method, ct string) *httpTransport {
	if ct == "" {
		ct = "application/drpc"
	}
	return &httpTransport{c: c, method: method, endpoint: endpoint, ct: ct}
}

func (t *httpTransport) Read(buf []byte) (int, error) {
	req, err := http.NewRequest(t.method, t.endpoint, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", t.ct)
	rsp, err := t.c.Do(req)
	if err != nil {
		return 0, err
	}
	buf, err = ioutil.ReadAll(rsp.Body)
	return len(buf), err
}

func (t *httpTransport) Write(buf []byte) (int, error) {
	req, err := http.NewRequest(t.method, t.endpoint, bytes.NewReader(buf))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", t.ct)
	_, err = t.c.Do(req)
	if err != nil {
		return 0, err
	}
	return len(buf), err
}

func (t *httpTransport) Close() error {
	t.closed = true
	return nil
}
