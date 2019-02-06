package gurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func NewRenderer() *Renderer {
	return &Renderer{
		Status: DefaultStatusRender,
		Header: DefaultHeaderRender,
		Plain:  PlainRender,
		JSON:   JSONRender,
		XML:    XMLRender,
	}
}

func (r *Renderer) render(res *http.Response) error {
	if r.Status != nil {
		fmt.Print(r.Status(res.Status, res.StatusCode))
	}

	if r.Header != nil {
		fmt.Print(r.Header(res.Header))
	}

	var bodyRender BodyRender
	cType := res.Header.Get("Content-Type")
	if strings.Contains(cType, "application/json") {
		bodyRender = r.JSON
	} else if strings.Contains(cType, "application/xml") {
		bodyRender = r.XML
	} else {
		bodyRender = r.Plain
	}

	if bodyRender != nil {
		defer res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		fmt.Print(bodyRender(string(bodyBytes)))
	}
	return nil
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

var sectionColor = color.New(color.Underline)

func sectionStr(str string) string {
	return fmt.Sprintf("%s\n", sectionColor.Sprint(str))
}

func statusStr(status string, code int) string {
	return color.New(colorForStatus(code)).Sprintf("%s", status)
}

func DefaultStatusRender(status string, code int) string {
	return fmt.Sprintf("%s%s\n\n", sectionStr("> Status"), statusStr(status, code))
}

func DefaultHeaderRender(h http.Header) string {
	var b bytes.Buffer
	b.WriteString(sectionStr("> Header"))
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
	return fmt.Sprintf("%s%s\n\n", sectionStr("> Body"), body)
}

func JSONRender(body string) string {
	var b bytes.Buffer
	if err := json.Indent(&b, []byte(body), "", "  "); err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s\n\n", sectionStr("> Body"), b.String())
}

func XMLRender(body string) string {
	x := xmlfmt.FormatXML(body, "", "  ")
	return fmt.Sprintf("%s%s\n\n", sectionStr("> Body"), x)
}
