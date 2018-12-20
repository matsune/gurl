package gurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func renderHeader(h http.Header) {
	fmt.Println("[Header]")
	for k, arr := range h {
		str := fmt.Sprintf("%s: ", k)
		for i, v := range arr {
			if i != 0 {
				str += fmt.Sprintf(", ")
			}
			str += fmt.Sprintf("%s", v)
		}
		fmt.Println(str)
	}
	fmt.Println()
}

type BodyRender func(body []byte) (string, error)

func plainRender(body []byte) (string, error) {
	return string(body), nil
}

func jsonRender(body []byte) (string, error) {
	var b bytes.Buffer
	if err := json.Indent(&b, body, "", "\t"); err != nil {
		return "", err
	}
	return b.String(), nil
}
