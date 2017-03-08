package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cbroglie/mustache"
	"github.com/hashicorp/terraform/helper/variables"
)

func parseCmdArgs() (string, map[string]interface{}, error) {
	vars := map[string]interface{}{}
	var file string
	flag.Var((*variables.Flag)(&vars), "var", "variables")
	flag.Var((*variables.FlagFile)(&vars), "var-file", "variable file")
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

func renderTemplate(template string, vars map[string]interface{}) (string, error) {
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
