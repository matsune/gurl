package gurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgsCmdName(t *testing.T) {
	args := []string{"./gurl"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Equal(t, args[0], cmdArgs.cmdName)
		assert.Empty(t, cmdArgs.rest)
	}
}

func TestParseArgsVersion(t *testing.T) {
	args := []string{"./gurl", "-v"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.True(t, cmdArgs.flags.Version)
	}

	args = []string{"./gurl", "--version"}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.True(t, cmdArgs.flags.Version)
	}
}

func TestParseArgsInteractive(t *testing.T) {
	args := []string{"./gurl", "-i"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.True(t, cmdArgs.isInteractive)
	}

	args = []string{"./gurl", "--interactive"}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.True(t, cmdArgs.isInteractive)
	}
}

func TestParseArgsOutOneline(t *testing.T) {
	args := []string{"./gurl", "-o"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.True(t, cmdArgs.flags.OutOneline)
	}

	args = []string{"./gurl", "--out-oneline"}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.True(t, cmdArgs.flags.OutOneline)
	}
}

func TestParseArgsBasic(t *testing.T) {
	args := []string{"gurl", "-u=user"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.Equal(t, "user", cmdArgs.flags.Basic)
	}

	args = []string{"gurl", "-u=user:pass"}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.Equal(t, "user:pass", cmdArgs.flags.Basic)
	}

	args = []string{"gurl", "--user=user:pass"}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.Equal(t, "user:pass", cmdArgs.flags.Basic)
	}
}

func TestParseArgsHeader(t *testing.T) {
	args := []string{"gurl", "-H", "A:B", "-H", "c:d"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.ElementsMatch(t, []string{"A:B", "c:d"}, cmdArgs.flags.Headers)
	}

	args = []string{"gurl", "--header=A:B", `--header="'c:':'d:'"`}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.ElementsMatch(t, []string{"A:B", `'c:':'d:'`}, cmdArgs.flags.Headers)
	}
}

func TestParseArgsJSON(t *testing.T) {
	args := []string{"gurl", "-j", `{"user": "u", "password": "p"}`}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.NotNil(t, cmdArgs.flags.JSON)
		assert.Equal(t, `{"user": "u", "password": "p"}`, cmdArgs.flags.JSON)
	}

	args = []string{"gurl", "--json", `{"user": "u", "password": "p"}`}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.NotNil(t, cmdArgs.flags.JSON)
		assert.Equal(t, `{"user": "u", "password": "p"}`, cmdArgs.flags.JSON)
	}
}

func TestParseArgsXML(t *testing.T) {
	args := []string{"gurl", "-x", `<user>u</user><password>p</password>`}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.NotNil(t, cmdArgs.flags.XML)
		assert.Equal(t, `<user>u</user><password>p</password>`, cmdArgs.flags.XML)
	}

	args = []string{"gurl", "--xml", `<user>u</user><password>p</password>`}
	cmdArgs, err = parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.NotNil(t, cmdArgs.flags.XML)
		assert.Equal(t, `<user>u</user><password>p</password>`, cmdArgs.flags.XML)
	}
}

func TestParseArgsForm(t *testing.T) {
	args := []string{"gurl", "-f", "user:u", "--form=password:p"}
	cmdArgs, err := parseArgs(args)
	if assert.NoError(t, err) {
		assert.Empty(t, cmdArgs.rest)
		assert.ElementsMatch(t, []string{"user:u", "password:p"}, cmdArgs.flags.Form)
	}
}

// Fail cases

func TestParseArgsEmpty(t *testing.T) {
	args := []string{}
	_, err := parseArgs(args)
	assert.Error(t, err)
}

func TestParseArgsBasicFail(t *testing.T) {
	args := []string{"gurl", "--user"}
	_, err := parseArgs(args)
	assert.Error(t, err)
}
func TestParseArgsHeaderFail(t *testing.T) {
	args := []string{"gurl", "-H"}
	_, err := parseArgs(args)
	assert.Error(t, err)
}

func TestParseArgsJSONFail(t *testing.T) {
	args := []string{"gurl", "--json"}
	_, err := parseArgs(args)
	assert.Error(t, err)
}

func TestParseArgsXMLFail(t *testing.T) {
	args := []string{"gurl", "--xml"}
	_, err := parseArgs(args)
	assert.Error(t, err)
}

func TestParseArgsEncodedFail(t *testing.T) {
	args := []string{"gurl", "--data"}
	_, err := parseArgs(args)
	assert.Error(t, err)
}
