package main

import (
	"log"
)

type Admin struct {
	rules    []AclRule
	consumer string
	method   string
}

func (adm Admin) Logging(nothing *Nothing, admLs Admin_LoggingServer) error {
	for {
		event := &Event{
			Timestamp: 0,
			Consumer:  adm.consumer,
			Method:    adm.method,
			Host:      "127.0.0.1:",
		}

		admLs.Send(event)
	}

	return nil
}

func (adm Admin) Statistics(statInterval *StatInterval, admSs Admin_StatisticsServer) error {
	log.Println("*Statistics()*")
	return nil
}
