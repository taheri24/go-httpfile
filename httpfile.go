package httpfile

import (
	"io"
	"net/http"
)

type HttpFile struct {
	requestMap map[string][]string
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
		requestMap: make(map[string][]string),
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
