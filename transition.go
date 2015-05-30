package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type Transition_List struct {
	Transitions []Transition
}

type Transition struct {
	Id   string
	Name string
	To   Status
}

func do_display_trans() {
	trans_list := do_get_trans()
	for _, element := range trans_list.Transitions {
		fmt.Printf("[%s] Name: %6s, To: %s\n", element.Id, element.Name, element.To.Name)
	}
}

func do_get_trans() Transition_List {
	Url, err := url.Parse(conf.Url)
	check(err)
	Url.Path += "/rest/api/2/issue/" + os.Getenv("current_issue") + "/transitions"
	parameters := url.Values{}
	parameters.Add("expand", "transitions.fields")
	Url.RawQuery = parameters.Encode()
	json_str := get_request(Url.String())
	var trans_list Transition_List
	err = json.Unmarshal(json_str, &trans_list)
	return trans_list
}
func do_trans_update(name string) {
	trans_list := do_get_trans()
	id := id_for_name(name, trans_list)
	Url, err := url.Parse(conf.Url)
	check(err)
	Url.Path += "/rest/api/2/issue/" + os.Getenv("current_issue") + "/transitions"
	data := "{ \"transition\": { \"id\": \"" + id + "\" } }"
	json_str := post_request(Url.String(), data)
	println(string(json_str))
}

func id_for_name(name string, trans Transition_List) string {
	for _, element := range trans.Transitions {
		if element.Name == name {
			return element.Id
		}
	}
	fmt.Printf("Transition %s not found aborting", name)
	os.Exit(1)
	return ""
}
