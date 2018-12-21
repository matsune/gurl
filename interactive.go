package gurl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"

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
		opts.SetBasic(u, p)
	}

	for {
		prompt := promptui.Select{
			Label: "Header",
			Items: []string{"End", "Add"},
			Templates: &promptui.SelectTemplates{
				Selected: fmt.Sprintf(``),
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return err
		}

		if idx == 0 {
			break
		}

		k, v, err := inputKeyValue()
		if err != nil {
			return nil
		}
		opts.SetHeader(k, v)
	}

	if opts.Body == nil {
		idx, err = selectItem("Data", []string{"None", "JSON", "XML", "Form"})
		if err != nil {
			return err
		}

		if idx == 1 {
			str, err := openEditor()
			if err != nil {
				return err
			}
			opts.Body = JSONData(str)
		} else if idx == 2 {
			str, err := openEditor()
			if err != nil {
				return err
			}
			opts.Body = XMLData(str)
		} else if idx == 3 {
			e := url.Values{}
			for {
				idx, err = selectItem("FormData", []string{"End", "Add"})
				if err != nil {
					return err
				}
				if idx == 0 {
					break
				}
				k, v, err := inputKeyValue()
				if err != nil {
					return nil
				}
				e.Set(k, v)
			}
			opts.Body = EncodedData(e)
		}
	}

	fmt.Println()

	g, err := New(opts)
	if err != nil {
		return err
	}

	if err := g.DoRequest(); err != nil {
		return err
	}

	return g.Render()
}

func openEditor() (string, error) {
	tmpFile := fmt.Sprintf("gurl.%s.tmp", randString(5))

	file, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	file.Close()

	cmd := exec.Command("vim", tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return "", err
	}

	file, err = os.OpenFile(tmpFile, os.O_RDONLY, 0600)
	if err != nil {
		return "", nil
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		file.Close()
		return "", nil
	}
	file.Close()
	if err = os.Remove(tmpFile); err != nil {
		return "", err
	}
	return string(bytes), nil
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
	idx, err := selectItem("Method", methods)
	return methods[idx], err
}

func inputURL() (string, error) {
	validate := func(input string) error {
		if !isURL(input) {
			return errors.New("not url")
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Success: fmt.Sprintf(`{{ "%s" | green }} {{ . }}: `, promptui.IconGood),
	}

	prompt := promptui.Prompt{
		Label:     "URL",
		Validate:  validate,
		Templates: templates,
	}
	return prompt.Run()
}

func selectItem(label string, items []string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Selected: fmt.Sprintf(`{{ "%s" | green }} %s: {{ . }}`, promptui.IconGood, label),
		},
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
		Templates: &promptui.PromptTemplates{
			Success: fmt.Sprint(`{{ . }}: `),
		},
	}
	return prompt.Run()
}

func inputPassword() (string, error) {
	prompt := promptui.Prompt{
		Label: "password",
		Mask:  '*',
		Templates: &promptui.PromptTemplates{
			Success: fmt.Sprint(`{{ . }}: `),
		},
	}
	return prompt.Run()
}

func inputKeyValue() (string, string, error) {
	prompt := promptui.Prompt{
		Label: "Key",
		Templates: &promptui.PromptTemplates{
			Success: fmt.Sprint(`{{ . }}: `),
		},
	}
	k, err := prompt.Run()
	if err != nil {
		return "", "", err
	}
	prompt = promptui.Prompt{
		Label: "Value",
		Templates: &promptui.PromptTemplates{
			Success: fmt.Sprint(`{{ . }}: `),
		},
	}
	v, err := prompt.Run()
	if err != nil {
		return "", "", err
	}
	return k, v, nil
}
