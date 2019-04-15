package gurl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/matsune/jc"
)

type BodyType int

const (
	BodyPlain BodyType = iota
	BodyJSON
	BodyXML
)

type (
	Renderer interface {
		Status(status string, code int) string
		Header(h http.Header) string
		Body(body string, ty BodyType) string
		Oneliner(str string) string
	}

	renderer struct {
		*jc.JC
	}
)

func NewRenderer() Renderer {
	j := jc.New()
	j.SetKeyColor(color.New(color.FgHiBlue))
	j.SetNumberColor(color.New(color.FgCyan))
	j.SetStringColor(color.New(color.FgGreen))
	j.SetBoolColor(color.New(color.FgHiRed))
	return &renderer{
		JC: j,
	}
}

func render(r Renderer, res *http.Response) error {
	if r == nil {
		return errors.New("Renderer is nil")
	}
	if res == nil {
		return errors.New("response is nil")
	}

	fmt.Print(r.Status(res.Status, res.StatusCode))

	fmt.Print(r.Header(res.Header))

	bt := BodyPlain
	cType := res.Header.Get("Content-Type")
	if strings.Contains(cType, "application/json") {
		bt = BodyJSON
	} else if strings.Contains(cType, "application/xml") {
		bt = BodyXML
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Print(r.Body(string(bodyBytes), bt))

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

func (r *renderer) section(str string) string {
	return fmt.Sprintf("%s", color.New(color.Underline).Sprint(str))
}

func (r *renderer) Status(status string, code int) string {
	return fmt.Sprintf("%s\n%s\n\n", r.section("> Status"), color.New(colorForStatus(code)).Sprintf("%s", status))
}

func (r *renderer) Header(h http.Header) string {
	var b bytes.Buffer
	b.WriteString(r.section("> Header\n"))
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

func (r *renderer) Body(body string, ty BodyType) string {
	fmt.Println(r.section("> Body"))
	switch ty {
	case BodyJSON:
		return r.json(body)
	case BodyXML:
		return r.xml(body)
	default:
		return fmt.Sprintf("%s\n\n", body)
	}
}

func (r *renderer) json(body string) string {
	var b bytes.Buffer
	r.SetWriter(&b)
	r.Colorize(body)
	return b.String()
}

func (r *renderer) xml(body string) string {
	x := xmlfmt.FormatXML(body, "", "  ")
	return fmt.Sprintf("%s\n\n", x)
}

func (r *renderer) Oneliner(str string) string {
	return fmt.Sprintf("%s\n%s\n", r.section("> one-liners"), str)
}
