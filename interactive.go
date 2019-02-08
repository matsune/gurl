package gurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey"
)

type interactivePrompt interface {
	askPassword(msg string) (string, error)
	confirm(msg string) (bool, error)
	selectOne(msg string, options []string) (string, error)
	input(msg string, validation func(interface{}) error) (string, error)
	inputEditor() (string, error)
}

type surveyPrompt struct {
	editor string
}

func (d *surveyPrompt) askPassword(msg string) (string, error) {
	p := ""
	s := &survey.Password{
		Message: msg,
	}
	if err := survey.AskOne(s, &p, nil); err != nil {
		return "", err
	}
	return p, nil
}

func (d *surveyPrompt) confirm(msg string) (res bool, err error) {
	prompt := &survey.Confirm{
		Message: msg,
	}
	err = survey.AskOne(prompt, &res, nil)
	return
}

func (d *surveyPrompt) selectOne(msg string, options []string) (res string, err error) {
	p := &survey.Select{
		Message: msg,
		Options: options,
	}
	err = survey.AskOne(p, &res, nil)
	return
}

func (d *surveyPrompt) input(msg string, validation func(interface{}) error) (res string, err error) {
	p := &survey.Input{
		Message: msg,
	}
	err = survey.AskOne(p, &res, validation)
	return
}

func (d *surveyPrompt) inputEditor() (string, error) {
	tmp, err := ioutil.TempFile("", "gurl")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command(d.editor, tmp.Name())
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

var prompt interactivePrompt

func init() {
	prompt = &surveyPrompt{
		editor: "vim",
	}
}

func isEmpty(str string) bool {
	return len(str) < 1
}

func goInteractive(opts *Options) error {

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
		m, err := prompt.selectOne("Choose method:", methods)
		if err != nil {
			return err
		}
		opts.Method = m
	}

	// URL

	if isEmpty(opts.URL) {
		validation := func(res interface{}) error {
			if _, ok := res.(string); !ok {
				return fmt.Errorf("not URL")
			}
			return nil
		}
		url, err := prompt.input("URL:", validation)
		if err != nil {
			return err
		}
		opts.URL = url
	}

	// Basic auth

	if opts.Basic == nil {
		res, err := prompt.confirm("Use Basic Authorization?")
		if err != nil {
			return err
		}
		if res {
			opts.Basic = &Basic{}
		}
	}

	if opts.Basic != nil {
		if isEmpty(opts.Basic.User) {
			user, err := prompt.input("User", nil)
			if err != nil {
				return err
			}
			opts.Basic.User = user
		}

		if isEmpty(opts.Basic.Password) {
			msg := fmt.Sprintf("Password for user %s:", opts.Basic.User)
			password, err := prompt.askPassword(msg)
			if err != nil {
				return err
			}
			opts.Basic.Password = password
		}
	}

	// Headers

	for {
		res, err := prompt.confirm("Add custom header?")
		if err != nil {
			return err
		}
		if !res {
			break
		}

		k, v, err := inputKeyValue()
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
		d, err := prompt.selectOne("Body:", options)
		if err != nil {
			return err
		}

		var body string
		switch d {
		case json:
			body, err = prompt.inputEditor()
			if err != nil {
				return err
			}
			opts.Body = JSONData(body)
		case xml:
			body, err = prompt.inputEditor()
			if err != nil {
				return err
			}
			opts.Body = XMLData(body)
		case form:
			urlValues, err := inputForm()
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

func inputForm() (url.Values, error) {
	e := url.Values{}

	k, v, err := inputKeyValue()
	if err != nil {
		return nil, err
	}
	e.Set(k, v)

	for {
		res, err := prompt.confirm("Add more form data?")
		if err != nil {
			return nil, err
		}
		if !res {
			break
		}

		k, v, err := inputKeyValue()
		if err != nil {
			return nil, err
		}
		e.Set(k, v)
	}
	return e, nil
}

func inputKeyValue() (k string, v string, err error) {
	k, err = prompt.input("Key:", nil)
	if err != nil {
		return
	}
	v, err = prompt.input("Value:", nil)
	if err != nil {
		return
	}
	return
}
