package main

import (
	"encoding/json"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"

	"github.com/GrooveCommunity/go-dispatcher/internal"
	svc "github.com/GrooveCommunity/go-dispatcher/service"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/validate-forward", handleValidateForward).Methods("POST")

	log.Println(os.Getenv("APP_PORT"))

	go internal.ForwardIssue(os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_TOKENAPI"), os.Getenv("JIRA_ENDPOINT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(svc.ValidateHealthy())
}

/*func handleForwardTickets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(svc.ForwardIssue(os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_TOKENAPI"), os.Getenv("JIRA_ENDPOINT")))
}*/

func handleValidateForward(w http.ResponseWriter, r *http.Request) {
	//json.NewEncoder(w).Encode(&svc.ValidateHealthy())
}
