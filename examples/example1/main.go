package main

import (
	"bytes"

	httpfile "github.com/taheri24/go-http"
)

func main() {
	buf := bytes.NewBufferString(`
 @base=222
# @name newSession
POST /ff
Content-Type: application/x-www-form-url

{"A":2
,"2":"A"}
	`)
	httpfile.New(buf)

}
