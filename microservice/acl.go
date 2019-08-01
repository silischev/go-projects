package main

import (
	"encoding/json"
)

type AclRule struct {
	route   string
	methods []string
}

func CreateRulesFromIncomingMessage(message []byte) ([]AclRule, error) {
	var aclIncomingMess map[string][]string
	err := json.Unmarshal(message, &aclIncomingMess)

	if err != nil {
		return nil, err
	}

	var rules []AclRule
	for route, methods := range aclIncomingMess {
		rules = append(rules, AclRule{route, methods})
	}

	return rules, nil
}
