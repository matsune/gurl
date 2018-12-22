package gurl

import (
	"encoding/base64"
	"fmt"
	"math/rand"
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

func split(str, sep string) (key, value string, err error) {
	kvs := strings.Split(str, sep)
	if len(kvs) != 2 {
		err = fmt.Errorf("invalid format '%s', please use 'key%svalue'.", str, sep)
		return
	}
	key = kvs[0]
	value = kvs[1]
	return
}

func splitKVs(kvs []string, sep string) (map[string][]string, error) {
	m := make(map[string][]string)
	for _, kv := range kvs {
		k, v, err := split(kv, sep)
		if err != nil {
			return nil, err
		}
		m[k] = append(m[k], v)
	}
	return m, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
