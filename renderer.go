package gurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-xmlfmt/xmlfmt"
)

type (
	StatusRender func(status string, code int) string
	HeaderRender func(h http.Header) string
	BodyRender   func(body string) string

	Renderer struct {
		Status StatusRender
		Header HeaderRender
		Plain  BodyRender
		JSON   BodyRender
		XML    BodyRender
	}
)

func DefaultRenderer() *Renderer {
	return &Renderer{
		Status: DefaultStatusRender,
		Header: DefaultHeaderRender,
		Plain:  PlainRender,
		JSON:   JSONRender,
		XML:    XMLRender,
	}
}

func colorForStatus(code int) color.Attribute {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return color.FgGreen
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return color.FgCyan
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return color.FgYellow
	default:
		return color.FgRed
	}
}

func DefaultStatusRender(status string, code int) string {
	d := color.New(colorForStatus(code))
	return fmt.Sprintf("Status: %s\n\n", d.Sprintf("%s", status))
}

const (
	headerSection = "--- Header ---\n\n"
	bodySection   = "---  Body  ---\n\n"
)

func DefaultHeaderRender(h http.Header) string {
	var b bytes.Buffer
	b.WriteString(headerSection)
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
	b.WriteString("\n")
	return b.String()
}

func PlainRender(body string) string {
	return fmt.Sprintf("%s%s", bodySection, body)
}

func JSONRender(body string) string {
	var b bytes.Buffer
	if err := json.Indent(&b, []byte(body), "", "  "); err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s", bodySection, b.String())
}

func XMLRender(body string) string {
	x := xmlfmt.FormatXML(body, "", "  ")
	return fmt.Sprintf("%s%s", bodySection, x)
}
