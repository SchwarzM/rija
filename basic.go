package main

import (
	"os/exec"
)

func get_request(url string) []byte {
	json_str, err := exec.Command(
		"/usr/bin/curl",
		"-s",
		"-E",
		conf.Cert+":"+conf.Pass,
		url,
	).Output()
	check(err)
	return json_str
}

func post_request(url string, data string) []byte {
	json_str, err := exec.Command(
		"/usr/bin/curl",
		"-i",
		"-X",
		"POST",
		"--data",
		data,
		"-H",
		"Content-Type: application/json",
		"-E",
		conf.Cert+":"+conf.Pass,
		url,
	).Output()
	check(err)
	println(string(json_str))
	return json_str
}
