package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	//"path/filepath"
	"net/url"
	"strconv"
)

var Commands = []cli.Command{
	print_issues,
	get_issues,
	set_issue,
}

var print_issues = cli.Command{
	Name:        "print",
	Usage:       "Prints the currently stored issues",
	Description: "Prints the currently stored issues",
	Action:      do_print_issues,
}

var get_issues = cli.Command{
	Name:        "get",
	Usage:       "Gets the list of currently known issues",
	Description: "Gets the list of currently known issues",
	Action:      do_get_issues,
}

var set_issue = cli.Command{
	Name:        "set",
	Usage:       "[n] Number of the Issue to set",
	Description: "Returns a sourceable string that sets the current_issue",
	Action:      do_set_issue,
}

var conf Conf

var list []Issue_List

type Issue_List struct {
	Key     string
	Summary string
}

type Issue_Fields struct {
	Summary string
}

type Issue struct {
	Key    string
	Fields Issue_Fields
}

type Response struct {
	Issues []Issue
}

func action(c *cli.Context) {
	println("boom i say")
}

func do_print_issues(c *cli.Context) {
	read_issues()
	for index, element := range list {
		fmt.Printf("[%d] %s: %s\n", index, element.Key, element.Summary)
	}
}

func do_get_issues(c *cli.Context) {
	get_configure()
	var resp Response
	Url, err := url.Parse(conf.Url)
	check(err)
	Url.Path += "/rest/api/2/search"
	parameters := url.Values{}
	parameters.Add("jql", "assignee="+conf.User+" AND ( status=Open OR status=\"In Progress\" )")
	Url.RawQuery = parameters.Encode()
	fmt.Printf("Enc Url %q\n", Url.String())
	json_str, err := exec.Command(
		"/usr/bin/curl",
		"-s",
		"-E",
		conf.Cert+":"+conf.Pass,
		Url.String(),
	).Output()
	check(err)
	err = json.Unmarshal(json_str, &resp)
	check(err)
	for _, element := range resp.Issues {
		item := Issue_List{element.Key, element.Fields.Summary}
		list = append(list, item)
	}
	write_issues()
}

func do_set_issue(c *cli.Context) {
	read_issues()
	if len(c.Args()) < 1 {
		println("Need a number of the issue to set")
	}
	index, err := strconv.Atoi(c.Args()[0])
	check(err)
	if index >= len(list) {
		fmt.Printf("No Issue at Index %d\n", index)
		os.Exit(1)
	}
	issue := list[index]
	fmt.Printf("set -gx current_issue %s\n", issue.Key)
}

type Conf struct {
	User string
	Pass string
	Url  string
	Cert string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func read_issues() {
	//filename, _ := filepath.Abs("~/.config/rija/rija_issues.yml")
	yaml_str, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/rija/rija_issues.yml")
	check(err)
	err = yaml.Unmarshal(yaml_str, &list)
	check(err)
}

func write_issues() {
	//	filename, _ := filepath.Abs("./rija_issues.yml")
	filename := os.Getenv("HOME") + "/.config/rija/rija_issues.yml"
	yaml, err := yaml.Marshal(&list)
	check(err)
	err = ioutil.WriteFile(filename, yaml, 0600)
	check(err)

}

func get_configure() {
	config, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/rija/rija.yml")
	//	filename, _ := filepath.Abs("~/.config/rija/rija.yml")
	//config, err := ioutil.ReadFile(filename)
	check(err)
	err = yaml.Unmarshal(config, &conf)
	check(err)
	//	fmt.Printf("%v\n", conf)
}
