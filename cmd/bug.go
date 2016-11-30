package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/containous/flaeg"
	"net/url"
	"os/exec"
	"runtime"
	"text/template"
)

var (
	bugtracker  = "https://github.com/containous/traefik/issues/new"
	bugTemplate = `Please answer these questions before submitting your issue. Thanks!

### What version of Traefik are you using?
` + "```" + `
{{.Version}}
` + "```" + `

### What is your environment & configuration?
` + "```" + `
{{.Configuration}}
` + "```" + `

### What did you do?


### What did you expect to see?


### What did you see instead?
`
)

// NewBugCmd builds a new Bug command
func NewBugCmd(getConf func() interface{}) *flaeg.Command {

	//version Command init
	return &flaeg.Command{
		Name:                  "bug",
		Description:           `Report an issue on Traefik bugtracker`,
		Config:                struct{}{},
		DefaultPointersConfig: struct{}{},
		Run: func() error {
			var version bytes.Buffer
			if err := getVersionPrint(&version); err != nil {
				return err
			}

			tmpl, err := template.New("").Parse(bugTemplate)
			if err != nil {
				return err
			}

			configJSON, err := json.MarshalIndent(getConf(), "", " ")
			if err != nil {
				return err
			}

			v := struct {
				Version       string
				Configuration string
			}{
				Version:       version.String(),
				Configuration: string(configJSON),
			}

			var bug bytes.Buffer
			if err := tmpl.Execute(&bug, v); err != nil {
				return err
			}

			body := bug.String()
			url := bugtracker + "?body=" + url.QueryEscape(body)
			if err := openBrowser(url); err != nil {
				fmt.Print("Please file a new issue at " + bugtracker + " using this template:\n\n")
				fmt.Print(body)
			} else {
				fmt.Print("Opening issue with:\n\n")
				fmt.Print(body)
			}

			return nil
		},
	}
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
