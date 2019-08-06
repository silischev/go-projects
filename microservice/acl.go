package main

import (
	"encoding/json"
	"regexp"
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

func hasAccess(consumer string, method string, rules []AclRule) bool {
	for _, rule := range rules {
		if rule.user == consumer {
			for _, aclMethod := range rule.methods {
				matched, _ := regexp.Match(aclMethod, []byte(method))

				if matched {
					return true
				}
			}
		}
	}

	return false
}
