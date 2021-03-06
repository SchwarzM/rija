package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	//"path/filepath"
	"net/url"
	"strconv"
	"strings"
)

var Commands = []cli.Command{
	print_issues,
	get_issues,
	set_issue,
	comment_issue,
	transition_issue,
}

var transition_issue = cli.Command{
	Name:        "transition",
	Usage:       "Transition the current issue",
	Description: "Displays possible transitions or transitions the issue",
	Action:      do_trans_issue,
}

var comment_issue = cli.Command{
	Name:        "comment",
	Usage:       "Comment on the current_issue",
	Description: "Comment on the current_issue",
	Action:      do_comment_issue,
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
	Status  string
}

type Status struct {
	Name string
}

type Issue_Fields struct {
	Summary string
	Status  Status
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
		fmt.Printf("[%2d] %s: <%11s> %s\n", index, element.Key, element.Status, element.Summary)
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
	json_str := get_request(Url.String())
	err = json.Unmarshal(json_str, &resp)
	check(err)
	for _, element := range resp.Issues {
		item := Issue_List{element.Key, element.Fields.Summary, element.Fields.Status.Name}
		list = append(list, item)
	}
	write_issues()
}

func do_set_issue(c *cli.Context) {
	read_issues()
	if len(c.Args()) < 1 {
		println("Need a number of the issue to set")
		os.Exit(1)
	}
	index, err := strconv.Atoi(c.Args()[0])
	check(err)
	if index >= len(list) {
		fmt.Printf("No Issue at Index %d\n", index)
		os.Exit(1)
	}
	issue := list[index]
	if strings.Contains(os.Getenv("SHELL"), "fish") {
		fmt.Printf("set -Ux current_issue %s\n", issue.Key)
	} else {
		fmt.Printf("export current_issue=%s\n", issue.Key)
	}
}

func do_trans_issue(c *cli.Context) {
	get_configure()
	if len(c.Args()) < 1 {
		do_display_trans()
	} else {
		do_trans_update(c.Args()[0])
	}
}

func do_comment_issue(c *cli.Context) {
	get_configure()
	if len(c.Args()) < 1 {
		println("Need a comment message to add")
		os.Exit(1)
	}
	message := c.Args()[0]
	println(message)
	Url, err := url.Parse(conf.Url)
	check(err)
	Url.Path += "/rest/api/2/issue/" + os.Getenv("current_issue") + "/comment"
	body := "{ \"body\": \"" + message + "\" }"
	post_request(Url.String(), body)
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
