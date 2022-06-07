package httpfile

import (
	"io"
	"net/http"
)

type HttpFile struct {
	requestMap map[string]*CapturedHttpRequest
}

type CapturedHttpRequest struct {
	reqName    string
	httpMethod string
	path       string
	headers    http.Header
	payload    []string
}

func New(contents io.Reader) *HttpFile {
	httpFile := &HttpFile{
		requestMap: make(map[string]*CapturedHttpRequest),
	}
	httpFile.Load(contents)
	return httpFile
}

type LineState int

const (
	BlankLine LineState = iota
	ReadingHeader
	ReadingPayload
)

var (
	httpMethods = []string{"GET", "HEAD", "POST"}
)

func (f *HttpFile) GetRequestMap() map[string]*CapturedHttpRequest {
	return f.requestMap
}
