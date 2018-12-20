package gurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Gurl struct {
	req *http.Request
	res *http.Response
}

func New(opts *Options) (*Gurl, error) {
	req, err := buildRequest(opts)
	if err != nil {
		return nil, err
	}
	return &Gurl{
		req: req,
		res: nil,
	}, nil
}

func (g *Gurl) DoRequest() error {
	if g.req == nil {
		return fmt.Errorf("request is nil")
	}

	c := new(http.Client)
	res, err := c.Do(g.req)
	if err != nil {
		return err
	}
	g.res = res
	return nil
}

func (g *Gurl) Render() error {
	renderHeader(g.res.Header)

	var renderFunc BodyRender
	if strings.Contains(g.res.Header.Get("Content-Type"), "application/json") {
		renderFunc = jsonRender
	} else {
		renderFunc = plainRender
	}

	defer g.res.Body.Close()
	body, err := ioutil.ReadAll(g.res.Body)
	if err != nil {
		return err
	}
	str, err := renderFunc(body)
	if err != nil {
		return err
	}
	fmt.Println("[Body]")
	fmt.Println(str)
	return nil
}
