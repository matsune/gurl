package gurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func (a *App) interactive(opts *Options) error {

	// method

	if isEmpty(opts.Method) {
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
		m, err := a.SelectOne("Select method:", methods)
		if err != nil {
			return err
		}
		opts.Method = m
	}

	// URL

	if isEmpty(opts.URL) {
		url, err := a.InputText("URL:")
		if err != nil {
			return err
		}
		opts.URL = url
	}

	// Basic auth

	if opts.Basic == nil {
		res, err := a.Confirm("Use Basic Authorization?")
		if err != nil {
			return err
		}
		if res {
			opts.Basic = &Basic{}
		}
	}

	if opts.Basic != nil {
		if isEmpty(opts.Basic.User) {
			user, err := a.InputText("User")
			if err != nil {
				return err
			}
			opts.Basic.User = user
		}

		if isEmpty(opts.Basic.Password) {
			msg := fmt.Sprintf("Password for user %s:", opts.Basic.User)
			password, err := a.InputPassword(msg)
			if err != nil {
				return err
			}
			opts.Basic.Password = password
		}
	}

	// Headers

	for {
		res, err := a.Confirm("Add custom header?")
		if err != nil {
			return err
		}
		if !res {
			break
		}

		k, v, err := a.inputKeyValue()
		if err != nil {
			return err
		}
		opts.CustomHeader.Add(k, v)
	}

	// Body

	if opts.Body == nil {
		const (
			none = "None"
			json = "JSON"
			xml  = "XML"
			form = "Form"
		)
		options := []string{none, json, xml, form}
		d, err := a.SelectOne("Body:", options)
		if err != nil {
			return err
		}

		var body string
		switch d {
		case json:
			body, err = a.inputEditor()
			if err != nil {
				return err
			}
			opts.Body = JSONData(body)
		case xml:
			body, err = a.inputEditor()
			if err != nil {
				return err
			}
			opts.Body = XMLData(body)
		case form:
			urlValues, err := a.inputForm()
			if err != nil {
				return err
			}
			opts.Body = EncodedData(urlValues)
		default:
			break
		}

	}

	fmt.Println()

	return nil
}

func (a *App) inputForm() (url.Values, error) {
	e := url.Values{}

	k, v, err := a.inputKeyValue()
	if err != nil {
		return nil, err
	}
	e.Set(k, v)

	for {
		res, err := a.Confirm("Add more form data?")
		if err != nil {
			return nil, err
		}
		if !res {
			break
		}

		k, v, err := a.inputKeyValue()
		if err != nil {
			return nil, err
		}
		e.Set(k, v)
	}
	return e, nil
}

func (a *App) inputKeyValue() (k string, v string, err error) {
	k, err = a.InputText("Key:")
	if err != nil {
		return
	}
	v, err = a.InputText("Value:")
	if err != nil {
		return
	}
	return
}

func (a *App) inputEditor() (string, error) {
	tmp, err := ioutil.TempFile("", "gurl")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command(a.Editor(), tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return "", err
	}

	tmp, err = os.OpenFile(tmp.Name(), os.O_RDONLY, 0600)
	if err != nil {
		return "", nil
	}
	defer os.Remove(tmp.Name())

	bytes, err := ioutil.ReadAll(tmp)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}
