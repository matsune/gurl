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

const editor = "vim"

func runInteractive(opts *Options) error {

	if err := selectMethod(opts); err != nil {
		return err
	}

	if err := inputURL(opts); err != nil {
		return err
	}

	if err := inputBasic(opts); err != nil {
		return err
	}

	if err := inputHeaders(opts); err != nil {
		return err
	}

	if err := inputBody(opts); err != nil {
		return err
	}

	fmt.Println()

	return nil
}

func selectMethod(opts *Options) error {
	if len(opts.Method) > 0 {
		return nil
	}
	m, err := _selectMethod()
	if err != nil {
		return err
	}
	opts.Method = m
	return nil
}

func _selectMethod() (string, error) {
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
	p := &survey.Select{
		Message: "Choose method:",
		Options: methods,
	}
	var res string
	err := survey.AskOne(p, &res, nil)
	return res, err
}

func inputURL(opts *Options) error {
	if len(opts.URL) > 0 {
		return nil
	}

	url, err := _inputURL()
	if err != nil {
		return err
	}
	opts.URL = url
	return nil
}

func _inputURL() (string, error) {
	var res string
	p := &survey.Input{
		Message: "URL:",
	}
	v := func(res interface{}) error {
		if v, ok := res.(string); !ok || !isURL(v) {
			return fmt.Errorf("not URL")
		}
		return nil
	}
	err := survey.AskOne(p, &res, v)
	return res, err
}

func inputBasic(opts *Options) error {
	if opts.Basic != nil && len(opts.Basic.User) > 0 && len(opts.Basic.Password) > 0 {
		return nil
	}

	if opts.Basic == nil {
		res := false
		prompt := &survey.Confirm{
			Message: "Use Basic Authorization?",
		}
		if err := survey.AskOne(prompt, &res, nil); err != nil {
			return err
		}
		if res {
			opts.Basic = &Basic{}
		} else {
			return nil
		}
	}

	if len(opts.Basic.User) == 0 {
		user := ""
		p := &survey.Input{
			Message: "User",
		}
		if err := survey.AskOne(p, &user, nil); err != nil {
			return err
		}
		opts.Basic.User = user
	}

	if len(opts.Basic.Password) == 0 {
		password := ""
		p := &survey.Password{
			Message: fmt.Sprintf("Password for user %s:", opts.Basic.User),
		}
		if err := survey.AskOne(p, &password, nil); err != nil {
			return err
		}
		opts.Basic.Password = password
	}
	return nil
}

func inputHeaders(opts *Options) error {
	for {
		res := false
		c := &survey.Confirm{
			Message: "Add custom header?",
		}
		if err := survey.AskOne(c, &res, nil); err != nil {
			return err
		}
		if !res {
			return nil
		}

		k, v, err := inputKeyValue()
		if err != nil {
			return err
		}
		opts.CustomHeader.Add(k, v)
	}
	return nil
}

func inputBody(opts *Options) error {
	if opts.Body != nil {
		return nil
	}

	const (
		none = "None"
		json = "JSON"
		xml  = "XML"
		form = "Form"
	)
	options := []string{none, json, xml, form}
	d := ""
	p := &survey.Select{
		Message: "Body:",
		Options: options,
	}
	if err := survey.AskOne(p, &d, nil); err != nil {
		return err
	}

	var err error
	switch d {
	case json:
		err = inputJSON(opts)
	case xml:
		err = inputXML(opts)
	case form:
		err = inputForm(opts)
	default:
		break
	}
	return err
}

func inputJSON(opts *Options) error {
	str, err := openEditor()
	if err != nil {
		return err
	}
	opts.Body = JSONData(str)
	return nil
}

func inputXML(opts *Options) error {
	str, err := openEditor()
	if err != nil {
		return err
	}
	opts.Body = XMLData(str)
	return nil
}

func inputForm(opts *Options) error {
	e := url.Values{}

	k, v, err := inputKeyValue()
	if err != nil {
		return err
	}
	e.Set(k, v)

	for {
		res := false
		c := &survey.Confirm{
			Message: "Add more form data?",
		}
		if err := survey.AskOne(c, &res, nil); err != nil {
			return err
		}

		if !res {
			break
		}

		k, v, err := inputKeyValue()
		if err != nil {
			return err
		}
		e.Set(k, v)
	}
	opts.Body = EncodedData(e)
	return nil
}

func inputKeyValue() (string, string, error) {
	k := ""
	p := &survey.Input{
		Message: "Key:",
	}
	if err := survey.AskOne(p, &k, nil); err != nil {
		return "", "", err
	}

	v := ""
	p = &survey.Input{
		Message: "Value:",
	}
	if err := survey.AskOne(p, &v, nil); err != nil {
		return "", "", err
	}
	return k, v, nil
}

func openEditor() (string, error) {
	tmp, err := ioutil.TempFile("", "gurl")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command(editor, tmp.Name())
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
