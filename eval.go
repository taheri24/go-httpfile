package httpfile

import "strings"

func (hf *HttpFile) Eval(s string) string {
	tmp := s + "     "
	var b strings.Builder
	length := len(s)
	jump := func(substr string, startFrom int) int {
		x := len(substr)

		for i := startFrom; i < length; i++ {
			sx := tmp[i : i+x]
			if sx == substr {
				return i - startFrom
			}
		}
		panic(`httpfile syntax broken`)
	}
	compute := func(varKey string) []byte {
		return []byte(varKey)
	}
	for i := 0; i < length; {
		s1 := tmp[i : i+1]
		s2 := tmp[i : i+2]
		s3 := tmp[i : i+3]

		if s3 == "\"{{" {
			jumpChCount := jump("}}\"", i)
			subExpr := s[i+2 : i+jumpChCount]
			b.Write(compute(subExpr))
			i += jumpChCount + 3

		} else if s2 == "{{" {
			jumpChCount := jump("}}", i)
			subExpr := s[i+2 : i+jumpChCount]
			b.Write(compute(subExpr))
			i += jumpChCount + 2
		} else {
			b.WriteString(s1)
			i++
		}
	}
	return b.String()
}
