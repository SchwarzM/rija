# rija

Commandline tool for getting jira issues assigned to me that are open oder "in progress"

## Usage

``` rija get ``` - gets the currently open issues from jira  
``` rija print ``` - print the list of currently known jira issues  
``` rija set <n> ``` - source the output to set the current_issue variable to the key of the choosen issue 

## Installation

- ```go get github.com/schwarzm/rija```  
- ```mkdir -p ~/.config/rija```
- ```cp example.yml ~/.config/rija/```
- Change the yml to your values

