package gurl

import "net/url"

type (
	BodyData interface {
		ContentType() string
		Raw() string
	}

	JSONData    string
	XMLData     string
	EncodedData url.Values
)

func (JSONData) ContentType() string {
	return "application/json"
}
func (XMLData) ContentType() string {
	return "application/xml"
}
func (EncodedData) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (j JSONData) Raw() string {
	return string(j)
}
func (x XMLData) Raw() string {
	return string(x)
}
func (f EncodedData) Raw() string {
	return url.Values(f).Encode()
}
