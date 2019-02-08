package gurl

import (
	"github.com/AlecAivazis/survey"
)

// TODO: mock test
type Prompt interface {
	InputText(msg string) (string, error)
	InputPassword(msg string) (string, error)
	Confirm(msg string) (bool, error)
	SelectOne(msg string, options []string) (string, error)
	Editor() string
}

type surveyPrompt struct {
	editor string
}

func (d *surveyPrompt) InputText(msg string) (res string, err error) {
	p := &survey.Input{
		Message: msg,
	}
	err = survey.AskOne(p, &res, nil)
	return
}

func (d *surveyPrompt) InputPassword(msg string) (string, error) {
	p := ""
	s := &survey.Password{
		Message: msg,
	}
	if err := survey.AskOne(s, &p, nil); err != nil {
		return "", err
	}
	return p, nil
}

func (d *surveyPrompt) Confirm(msg string) (res bool, err error) {
	prompt := &survey.Confirm{
		Message: msg,
	}
	err = survey.AskOne(prompt, &res, nil)
	return
}

func (d *surveyPrompt) SelectOne(msg string, options []string) (res string, err error) {
	p := &survey.Select{
		Message: msg,
		Options: options,
	}
	err = survey.AskOne(p, &res, nil)
	return
}

func (d *surveyPrompt) Editor() string {
	return d.editor
}
