package httpfile_test

import (
	"bytes"
	"testing"

	httpfile "github.com/taheri24/go-http"
)

func Test(t *testing.T) {
	s := bytes.NewBufferString(`# @name test22
POST /testing
Content-Type: application/json

{"dd":2}`)
	hf1 := httpfile.New(s)
	dd := hf1.GetRequestMap()["test22"]
	t.Logf("Test %v", dd)
	t.Fail()
}

func Test2(t *testing.T) {
	s := bytes.NewBufferString(``)
	hf1 := httpfile.New(s)
	dd := hf1.Eval(`{{baseURL}}/sdf`)
	t.Logf("Test %v", dd)
	t.Fail()
}
