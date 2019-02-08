package gurl

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// ignore upper and lower case
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

func splitKV(str string) (key, value string, err error) {
	kvs := strings.Split(str, ":")
	if len(kvs) != 2 {
		err = fmt.Errorf("invalid key:value format: %s", str)
		return
	}
	key = kvs[0]
	value = kvs[1]
	return
}

func splitKVs(kvs []string) (map[string][]string, error) {
	m := map[string][]string{}
	for _, kv := range kvs {
		k, v, err := splitKV(kv)
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

func isEmpty(str string) bool {
	return len(str) < 1
}
