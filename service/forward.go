package service

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
)

type Rule struct {
	Name    string `json:"name,omitempty"`
	Field   string `json:"field,omitempty"`
	Value   string `json:"value,omitempty"`
	Content string `json:"content,omitempty"`
}

type Issue struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Reporter    string `json:"reporter,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
}

type Response struct {
	Issues []Issue `json:"issues,omitempty"`
}

func ForwardIssue(username, token, endpoint string) Response {
	tp := jira.BasicAuthTransport{
		Username: username, //usuário do jira
		Password: token,    //token de api
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(endpoint))
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return Response{}
	}

	jql := "project = 'service desk' and type = incidente and status = 'AGUARDANDO SD'"

	rule := Rule{Name: "RulePortalClienteTEFComAnexo",
		Field: "Produtos ServiceDesk", Value: "Portal Cliente (TEF)", Content: "reexportação",
	}

	jql = getJql(rule, jql)

	issuesJira, err := getAllIssues(client, jql)

	if err != nil {
		panic(err)
	}

	var issues []Issue

	for _, v := range issuesJira {
		issues = append(issues, Issue{ID: v.ID, Description: v.Fields.Description, Reporter: v.Fields.Reporter.Name, Assignee: v.Fields.Assignee.Name})
	}

	return Response{Issues: issues}
}

func getAllIssues(client *jira.Client, searchString string) ([]jira.Issue, error) {
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

func getJql(rule Rule, jql string) string {
	return jql + " '" + rule.Field + "' = '" + rule.Field + "' and text ~ '" + rule.Content

}
