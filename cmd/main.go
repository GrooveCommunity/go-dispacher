package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GrooveCommunity/go-dispatcher/entity"

	"os"

	"github.com/gorilla/mux"

	"github.com/GrooveCommunity/go-dispatcher/internal"
)

var rules []entity.Rule

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/put-rule", handlePutRule).Methods("POST")

	log.Println(os.Getenv("APP_PORT"))

	rules = internal.GetRules()

	go internal.ForwardIssue(rules, os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_TOKENAPI"), os.Getenv("JIRA_ENDPOINT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func GetRules()[]entity.Rule {
     return rules	
}

func handlePutRule(w http.ResponseWriter, r *http.Request) {
	var rule entity.Rule

	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		panic(err)
	}

	internal.WriteRule(rule)
	
	log.Println("Regra escrita", rule)

	rules = append(rules, rule)
}
