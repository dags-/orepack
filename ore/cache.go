package ore

import (
	"bufio"
	"bytes"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	expire           = time.Second * 30
	update           = time.Second * 30
	getResponseCache = cache.New(expire, update)
)

func HttpGet(url string) (*http.Response, error) {
	rq, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}
	r, e := computeIfAbsent(rq.URL.String())
	if e != nil {
		return nil, e
	}
	return http.ReadResponse(bufio.NewReader(r), rq)
}

func computeIfAbsent(url string) (*bytes.Reader, error) {
	data, exists := getResponseCache.Get(url)
	if exists {
		return bytes.NewReader(data.([]byte)), nil
	}
	resp, e := doGet(url)
	if e != nil {
		return nil, e
	}
	getResponseCache.Set(url, resp, expire)
	return bytes.NewReader(resp), nil
}

func doGet(url string) ([]byte, error) {
	rs, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	buf := &bytes.Buffer{}
	rs.Write(buf)
	return buf.Bytes(), nil
}
