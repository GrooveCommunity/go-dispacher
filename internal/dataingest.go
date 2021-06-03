package internal

import (
	"log"

	"encoding/json"

	gcp "github.com/GrooveCommunity/glib-cloud-storage/gcp"
	"github.com/GrooveCommunity/go-dispatcher/entity"
)

func WriteRule(rule entity.Rule) {
	gcp.WriteObject(rule, "forward-dispatcher", rule.Name)
}

func GetRules() []entity.Rule {
	var rules []entity.Rule

	dataObjects := gcp.GetObjects("forward-dispatcher")

	for _, b := range dataObjects {
		var rule entity.Rule
		errUnmarsh := json.Unmarshal(b, &rule)

		if errUnmarsh != nil {
			log.Fatal("Erro no unmarshal\n", errUnmarsh.Error())
		}

		rules = append(rules, rule)
	}

	log.Println(rules)

	return rules
}
