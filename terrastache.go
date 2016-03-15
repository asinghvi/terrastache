package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cbroglie/mustache"
	"github.com/hashicorp/terraform/command"
)

func parseCmdArgs() (string, map[string]string, error) {
	vars := map[string]string{}
	var file string
	flag.Var((*command.FlagKV)(&vars), "var", "variables")
	flag.Var((*command.FlagKVFile)(&vars), "var-file", "variable file")
	flag.StringVar(&file, "template", "", "template file")
	flag.Parse()

	if file == "" {
		return "", nil, fmt.Errorf("must specify a template file using the template parameter")
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", nil, err
	}
	return string(b), vars, nil
}

func renderTemplate(template string, vars map[string]string) (string, error) {
	data, err := mustache.Render(template, vars)
	if err != nil {
		return "", err
	}
	return data, nil
}

func main() {
	template, vars, err := parseCmdArgs()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	rendered, err := renderTemplate(template, vars)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(rendered)
}
