package gurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Gurl struct {
	Req      *http.Request
	Res      *http.Response
	Renderer Renderer
}

func New(opts *Options) (*Gurl, error) {
	req, err := opts.buildRequest()
	if err != nil {
		return nil, err
	}

	return &Gurl{
		Req:      req,
		Res:      nil,
		Renderer: *NewRenderer(),
	}, nil
}

func (g *Gurl) DoRequest() error {
	if g.Req == nil {
		return fmt.Errorf("Gurl.req is nil")
	}

	c := new(http.Client)
	res, err := c.Do(g.Req)
	if err != nil {
		return err
	}
	g.Res = res
	return nil
}

func (g *Gurl) Render() error {
	if g.Renderer.Status != nil {
		fmt.Print(g.Renderer.Status(g.Res.Status, g.Res.StatusCode))
	}

	if g.Renderer.Header != nil {
		fmt.Print(g.Renderer.Header(g.Res.Header))
	}

	var bodyRender BodyRender
	cType := g.Res.Header.Get("Content-Type")
	if strings.Contains(cType, "application/json") {
		bodyRender = g.Renderer.JSON
	} else if strings.Contains(cType, "application/xml") {
		bodyRender = g.Renderer.XML
	} else {
		bodyRender = g.Renderer.Plain
	}

	if bodyRender != nil {
		defer g.Res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(g.Res.Body)
		if err != nil {
			return err
		}
		fmt.Print(bodyRender(string(bodyBytes)))
	}

	return nil
}
