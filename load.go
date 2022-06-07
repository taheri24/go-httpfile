package httpfile

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func parseRequestRouteLine(line string) (string, string) {
	for _, httpMethod := range httpMethods {
		prefix := fmt.Sprintf("%s ", httpMethod)
		if strings.HasPrefix(line, prefix) {
			segments := strings.SplitN(line, " ", 2)
			if len(segments) < 2 {
				return "", ""
			}
			return segments[0], segments[1]
		}
	}
	return "", ""
}
func trim(s string) string {
	return strings.Trim(s, "\r\n ")
}
func (hf *HttpFile) Load(contents io.Reader) {
	var (
		currentReq *CapturedHttpRequest
		lineState  LineState
	)
	captureCurrentRequest := func() {
		captured := false
		if currentReq != nil {
			if len(currentReq.reqName) > 0 && len(currentReq.httpMethod) > 0 && len(currentReq.path) > 0 {
				hf.requestMap[currentReq.reqName] = currentReq
				lineState = BlankLine
				captured = true
				fmt.Printf("capture a request \"%s\"  ", currentReq.reqName)
			}
		}
		if captured || currentReq == nil {
			currentReq = &CapturedHttpRequest{
				headers: make(http.Header),
				payload: []string{},
			}
		}
	}
	scanner := bufio.NewScanner(contents)
	captureCurrentRequest()
	for scanner.Scan() {
		line := scanner.Text()
		line = trim(line)
		if strings.HasPrefix(line, "##") {
			captureCurrentRequest()
		} else if lineState == ReadingPayload {
			currentReq.payload = append(currentReq.payload, line)
		} else if strings.HasPrefix(line, "@") {
			line = strings.Replace(line, "@", "", 1)
			segments := strings.Split(line, "=")
			fmt.Printf("detect variable %s=%s\n", segments[0], segments[1])
		} else if strings.HasPrefix(line, "#") && strings.Contains(line, "@name") {
			line = strings.Replace(line, "#", "", 1)
			line = strings.Replace(line, "@name", "", 1)
			fmt.Printf("detect request name %s \n", line)
			captureCurrentRequest()
			reqName := trim(line)
			currentReq.reqName = reqName
		} else if httpMethod, path := parseRequestRouteLine(line); httpMethod != "" {
			captureCurrentRequest()
			currentReq.httpMethod = httpMethod
			currentReq.path = path
			lineState = ReadingHeader
			fmt.Printf("detect request  [%s]  %s\n", httpMethod, path)
		} else if lineState == ReadingHeader && strings.Contains(line, ":") {
			segments := strings.Split(line, ":")
			headerKey, headerVal := trim(segments[0]), trim(segments[1])
			fmt.Printf("detect header request  %s  %s\n", headerKey, headerVal)
			currentReq.headers.Add(headerKey, headerVal)
		} else if line == "" && lineState == ReadingHeader && len(currentReq.httpMethod) > 0 && len(currentReq.path) > 0 {
			lineState = ReadingPayload
		}
	}
	captureCurrentRequest()
}
