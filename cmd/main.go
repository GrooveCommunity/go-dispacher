package main

import (
	"encoding/json"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"

	svc "github.com/GrooveCommunity/go-dispacher/service"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/forward-tickets", handleForwardTickets).Methods("POST")
	router.HandleFunc("/validate-forward", handleValidateForward).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(svc.ValidateHealthy())
}

func handleForwardTickets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(svc.ForwardIssue(os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_TOKENAPI"), os.Getenv("JIRA_ENDPOINT")))
}

func handleValidateForward(w http.ResponseWriter, r *http.Request) {
	//json.NewEncoder(w).Encode(&svc.ValidateHealthy())
}
