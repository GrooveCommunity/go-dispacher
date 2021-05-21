package main

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
)

func main() {
	tp := jira.BasicAuthTransport{
		Username: "", //usuário do jira
		Password: "", //token de api
	}

	jiraEndpoint := "" //url do jira

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(""))
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}

	//passar o nome do projeto, o tipo da issue e o status do card
	jql := "project = '' and type =  and status = ''"

	issues, err := GetAllIssues(client, jql)

	if err != nil {
		panic(err)
	}

	for _, v := range issues {
		fmt.Println(v.Fields.Description)
	}

}

func GetAllIssues(client *jira.Client, searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 10, // Max results can go up to 1000
			StartAt:    last,
		}

		chunk, resp, err := client.Issue.Search(searchString, opt)
		if err != nil {
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}
		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= total {
			return issues, nil
		}
	}

}
