package main

import (
	"log"
)

type Admin struct {
	rules []AclRule
}

func (adm Admin) Logging(nothing *Nothing, admLs Admin_LoggingServer) error {
	log.Println("*Logging()*")
	//event := &Event{Timestamp: 0, Consumer: "logger", Method: "/main.Admin/Logging", Host: ""}
	//admLs.Send(event)
	return nil
}

func (adm Admin) Statistics(statInterval *StatInterval, admSs Admin_StatisticsServer) error {
	log.Println("*Statistics()*")
	return nil
}
