package gurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-xmlfmt/xmlfmt"
)

type HeaderRender func(h http.Header) string

func DefaultHeaderRender(h http.Header) string {
	var b bytes.Buffer
	b.WriteString("[Header]\n")
	for k, arr := range h {
		b.WriteString(k + ": ")
		for i, v := range arr {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(v)
		}
		b.WriteString("\n")
	}
	return b.String()
}

type BodyRender func(body []byte) (string, error)

func PlainRender(body []byte) (string, error) {
	return fmt.Sprintf("[Body]\n%s", string(body)), nil
}

func JSONRender(body []byte) (string, error) {
	var b bytes.Buffer
	if err := json.Indent(&b, body, "", "  "); err != nil {
		return "", err
	}
	return fmt.Sprintf("[Body]\n%s", b.String()), nil
}

func XMLRender(body []byte) (string, error) {
	x := xmlfmt.FormatXML(string(body), "", "  ")
	return fmt.Sprintf("[Body]\n%s", x), nil
}
