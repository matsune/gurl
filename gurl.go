package gurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Gurl struct {
	Req          *http.Request
	Res          *http.Response
	HeaderRender HeaderRender
	BodyRender   BodyRender
}

func New(opts *Options) (*Gurl, error) {
	req, err := opts.buildRequest()
	if err != nil {
		return nil, err
	}

	return &Gurl{
		Req:          req,
		Res:          nil,
		HeaderRender: DefaultHeaderRender,
		BodyRender:   nil,
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
	fmt.Println(g.HeaderRender(g.Res.Header))

	bodyRender := g.BodyRender
	if g.BodyRender == nil {
		cType := g.Res.Header.Get("Content-Type")
		if strings.Contains(cType, "application/json") {
			bodyRender = JSONRender
		} else if strings.Contains(cType, "application/xml") {
			bodyRender = XMLRender
		} else {
			bodyRender = PlainRender
		}
	}

	defer g.Res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(g.Res.Body)
	if err != nil {
		return err
	}
	body, err := bodyRender(bodyBytes)
	if err != nil {
		return err
	}
	fmt.Println(body)
	return nil
}
