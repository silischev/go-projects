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

func hasAccess(consumer string, method string, rules []AclRule) (bool, error) {
	for _, rule := range rules {
		if rule.user == consumer {
			ok, err := isMethodAvailable(method, rule.methods)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}
		}
	}

	return false, nil
}

func isMethodAvailable(method string, rules []string) (bool, error) {
	for _, aclMethod := range rules {
		matched, err := regexp.Match(aclMethod, []byte(method))
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
