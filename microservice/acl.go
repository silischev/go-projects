package main

import (
	"encoding/json"
	"log"
)

type AclRule struct {
	user    string
	methods []string
}

func CreateRulesFromIncomingMessage(message []byte) ([]AclRule, error) {
	var aclIncomingMess map[string][]string
	err := json.Unmarshal(message, &aclIncomingMess)

	if err != nil {
		return nil, err
	}

	var rules []AclRule
	for user, methods := range aclIncomingMess {
		rules = append(rules, AclRule{user, methods})
	}

	return rules, nil
}

func hasAccess(context string, rules []AclRule) bool {
	for rule := range rules {
		//if
		log.Println(rule)
	}

	return false
}
