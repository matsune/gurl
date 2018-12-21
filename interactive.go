package gurl

import (
	"errors"
	"net/http"

	"github.com/manifoldco/promptui"
)

func runInteractive(opts *Options) error {
	if len(opts.Method) == 0 {
		m, err := selectMethod()
		if err != nil {
			return err
		}
		opts.Method = m
	}

	if len(opts.URL) == 0 {
		url, err := inputURL()
		if err != nil {
			return err
		}
		opts.URL = url
	}

	idx, err := selectItem("Authorization", []string{"None", "Basic Auth"})
	if err != nil {
		return err
	}
	if idx == 1 {
		u, err := inputUser()
		if err != nil {
			return err
		}
		p, err := inputPassword()
		if err != nil {
			return err
		}
		opts.setBasic(u, p)
	}

	g, err := New(opts)
	if err != nil {
		return err
	}

	if err := g.DoRequest(); err != nil {
		return err
	}

	return g.Render()
}

func selectMethod() (string, error) {
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
	idx, err := selectItem("Select Method", methods)
	return methods[idx], err
}

func inputURL() (string, error) {
	validate := func(input string) error {
		if !isURL(input) {
			return errors.New("not url")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "URL",
		Validate: validate,
	}
	return prompt.Run()
}

func selectItem(label string, items []string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return idx, nil
}

func inputUser() (string, error) {
	prompt := promptui.Prompt{
		Label: "user",
	}
	return prompt.Run()
}

func inputPassword() (string, error) {
	prompt := promptui.Prompt{
		Label: "password",
		Mask:  '*',
	}
	return prompt.Run()
}
