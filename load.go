package httpfile

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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
		reqLines       []string
		currentReqName string
		lineState      LineState
	)
	captureCurrentRequest := func() {
		captured := false
		if reqLines != nil && len(currentReqName) > 0 {
			hf.requestMap[currentReqName] = reqLines[:]
			lineState = BlankLine
			captured = true
			fmt.Printf("capture a request \"%s\"  ", currentReqName)
		}
		if captured || reqLines == nil {
			reqLines = make([]string, 0)
		}
	}
	scanner := bufio.NewScanner(contents)
	captureCurrentRequest()
	for scanner.Scan() {
		line := scanner.Text()
		line = trim(line)
		if strings.HasPrefix(line, "###") {
			captureCurrentRequest()
		} else if lineState == ReadingPayload {
			reqLines = append(reqLines, line)
		} else if strings.HasPrefix(line, "@") {
			line = strings.Replace(line, "@", "", 1)
			segments := strings.Split(line, "=")
			fmt.Printf("detect variable %s=%s\n", segments[0], segments[1])
		} else if strings.HasPrefix(line, "#") && strings.Contains(line, "@name") {
			line = strings.Replace(line, "#", "", 1)
			line = strings.Replace(line, "@name", "", 1)
			//fmt.Printf("detect request name %s \n", line)
			captureCurrentRequest()
			currentReqName = trim(line)

		} else if httpMethod, path := parseRequestRouteLine(line); httpMethod != "" {
			captureCurrentRequest()
			reqLines = append(reqLines, line)
			lineState = ReadingHeader
			fmt.Printf("detect request  [%s]  %s\n", httpMethod, path)
		} else if lineState == ReadingHeader && strings.Contains(line, ":") {
			reqLines = append(reqLines, line)
			segments := strings.Split(line, ":")
			headerKey, headerVal := trim(segments[0]), trim(segments[1])
			fmt.Printf("detect header request  %s  %s\n", headerKey, headerVal)
		} else if line == "" && lineState == ReadingHeader {
			lineState = ReadingPayload
			reqLines = append(reqLines, line)

		}
	}
	captureCurrentRequest()
}
func (hf *HttpFile) GetTemplate(name string) string {
	var lines []string
	if requestLines, found := hf.requestMap[name]; !found {
		lines = requestLines
	}
	sb := strings.Builder{}
	replacer := strings.NewReplacer("\"{{", "{{", "}}\"", "}}")
	var lineRegex *regexp.Regexp
	if temp, err := regexp.Compile(`"\{\{`); err != nil {
		lineRegex = temp
		panic(err)
	}
	//	payloadFormat := "json"
	for _, line := range lines {
		line = replacer.Replace(line)
		format := "json"
		line = lineRegex.ReplaceAllString(line, `{{.$1 | `+format+"}} ")
		//lineRegex.ReplaceAllString("")
	}
	return sb.String()
}
