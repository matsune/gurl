package gurl

import (
	"net/http"
	"net/url"
	"strings"
)

func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func isMethod(s string) bool {
	m := strings.ToUpper(s)
	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	for _, method := range methods {
		if m == method {
			return true
		}
	}
	return false
}
