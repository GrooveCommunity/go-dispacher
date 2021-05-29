package internal

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	//	"github.com/trivago/tgo/tcontainer"

	//	"reflect"

	"time"

	gcp "github.com/GrooveCommunity/glib-cloud-storage/gcp"
	//"github.com/fatih/structs"
)

/*type Rule struct {
	Name    string `json:"name,omitempty"`
	Field   string `json:"field,omitempty"`
	Value   string `json:"value,omitempty"`
	Content string `json:"content,omitempty"`
}*/

type Customfield10366 struct {
	Value string `json:"value"`
}

type Fields struct {
	Customfield10366 Customfield10366 `json:"customfield_10366"`
}

type DataField struct {
	Fields Fields `json:"fields"`
}

type Issue struct {
	ID                 string `json:"id,omitempty"`
	Description        string `json:"description,omitempty"`
	Reporter           string `json:"reporter,omitempty"`
	CreatedDate        string `json:"created_date,omitempty"`
	Type               string `json:"type,omitempty"`
	Priority           string `json:"priority,omitempty"`
	ProductServiceDesk string `json:"priority,omitempty"`
}

type Response struct {
	Issues []Issue `json:"issues,omitempty"`
}

func ForwardIssue(username, token, endpoint string) Response {

	for {

		fmt.Println("Buscando regras no GCS")
		dataObjects := gcp.GetObjects("forward-dispatcher")

		tp := jira.BasicAuthTransport{
			Username: username, //usuário do jira
			Password: token,    //token de api
		}

		client, err := jira.NewClient(tp.Client(), strings.TrimSpace(endpoint))
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			return Response{}
		}

		for _, dataObject := range dataObjects {
			jql := "project = 'service desk' and type = incidente and status = 'AGUARDANDO SD' and '" + dataObject.Forward.Name + "' = '" + dataObject.Forward.Value + "' and text ~ '" +
				dataObject.Forward.Content + "' and NOT attachments is EMPTY"

			issuesJira, err := getAllIssues(client, jql)

			if err != nil && !(strings.HasPrefix(err.Error(), "No response returned")) {
				fmt.Printf("\nError: %v\n", err)
				return Response{}
			}

			//		var issues []Issue

			//			customfield := tcontainer.NewMarshalMap()

			for _, v := range issuesJira {

				updateIssue(v.Key, "customfield_10366", "Squad PayReport")

				//issueService := client.Issue

				//cf, res, _ := issueService.GetCustomFields(v.ID)

				//issueJiraService, resIss, errIss := issueService.Get(v.ID, nil)

				//fmt.Println("GetIssue", issueJiraService, resIss, errIss)

				//m := structs.Map(v.Fields)

				//fmt.Println(v.Fields.Unknowns["customfield_10366"]["id"])
				//fmt.Println(m)

				//custom_field["customfield_10366"]

				//fmt.Println("Cfs: ", cf, res, errCf)

				//fmt.Println("Unknows: ", v.Fields.Unknowns)

				//var nmap tcontainer.MarshalMap
				//var p map[string]interface{}

				//nmap := v.Fields.Unknowns //.Value("customfield_10366")

				//nmap, _ := p.(tcontainer.MarshalMap)

				//mm, _ := (nmap.MarshalMap("customfield_10366"))
				//mmv, _ := mm.Value("value")

				//fmt.Println(reflect.TypeOf(mmv))

				//v.Fields.Unknowns["customfield_10366"] = dataObject.Forward.Squad
				//issue, _, err := client.Issue.Update(&v)

				/*if err != nil {
					fmt.Println("Erro na atualização da issue", err.Error())

					panic(err)
				}

				fmt.Printf("Issue %s encaminhada para o squad %s", issue.ID, dataObject.Forward.Squad)*/

				/*i := jira.Issue{

					Fields: &jira.IssueFields{
						Unknowns: customfield,
					},
				}

				issue, _, err := client.Issue.Update(&i)*/

				if err != nil {
					panic(err)
				}
			}

			//go DataIngest(issuesJira)

		}

		/*

			//rule := Rule{Name: "RulePortalClienteTEFComAnexo"} //Field: "Produtos ServiceDesk", Value: "Portal Cliente (TEF)", Content: "reexportação",

			//jql = getJql(rule, jql)

			issuesJira, err := getAllIssues(client, jql)

			if err != nil && !(strings.HasPrefix(err.Error(), "No response returned")) {
				fmt.Printf("\nError: %v\n", err)
				return Response{}
			}

			var issues []Issue

			for _, v := range issuesJira {
				createdDate, _ := v.Fields.Created.MarshalJSON()

				//log.Println(v.Fields.Unknowns)

				m := structs.Map(v.Fields)
				unknowns, okay := m["Unknowns"]

				if okay {
					for key, value := range unknowns.(tcontainer.MarshalMap) {
						//m[key] = value

						if key == "customfield_10519" {
							log.Println(value)
						}
					}
				}

				log.Println(m)

				//log.Println(v.Fields.Unknowns)

				issues = append(issues, Issue{ID: v.ID, Description: v.Fields.Description, Reporter: v.Fields.Reporter.DisplayName, CreatedDate: string(createdDate), Type: v.Fields.Type.Name, Priority: v.Fields.Priority.Name})
			}

			if err != nil && !(strings.HasPrefix(err.Error(), "No response returned")) {
				fmt.Printf("\nError: %v\n", err)
				return Response{}
			}

			//go DataIngest(issues)

			return Response{Issues: issues}

			/*

				if err != nil && !(strings.HasPrefix(err.Error(), "No response returned")) {
					fmt.Printf("\nError: %v\n", err)
					return Response{}
				}

				var issues []Issue

				for _, v := range issuesJira {

					createdDate, _ := v.Fields.Created.MarshalJSON()

					issues = append(issues, Issue{ID: v.ID, Description: v.Fields.Description, Reporter: v.Fields.Reporter.DisplayName, CreatedDate: string(createdDate), Type: v.Fields.Type.Name, Priority: v.Fields.Priority.Name})
				}

				go DataIngest(issues)

				return Response{Issues: issues}*/

		fmt.Println("Aguarando um minuto")

		time.Sleep(1 * time.Minute)
	}
}

func getAllIssues(client *jira.Client, searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
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

func updateIssue(keyID, customFieldID, customFieldValue string) {
	host := os.Getenv("JIRA_ENDPOINT") + "/rest/api/2/issue/" + keyID

	customfield10366 := Customfield10366{Value: customFieldValue}
	fields := Fields{Customfield10366: customfield10366}

	data := DataField{Fields: fields}
	jsonReq, err := json.Marshal(data)

	fmt.Println(string(jsonReq))

	req, err := http.NewRequest(http.MethodPut, host, bytes.NewBuffer(jsonReq))
	req.SetBasicAuth(os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_TOKENAPI"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
}
