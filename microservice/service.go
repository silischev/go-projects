package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type BizServerHandler interface {
	Check(context.Context, Nothing) Nothing
	Add(context.Context, Nothing) Nothing
	Test(context.Context, Nothing) Nothing
}

type Biz struct {
}

func (b Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	log.Println("In Check()")
	return &Nothing{Dummy: true}, nil
}

func (b Biz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	log.Println("In Add()")
	return &Nothing{Dummy: true}, nil
}

func (b Biz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	log.Println("In Test()")
	return &Nothing{Dummy: true}, nil
}

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	server := grpc.NewServer()
	biz := Biz{}

	var aclIncomingMess map[string]interface{}
	err := json.Unmarshal([]byte(ACLData), &aclIncomingMess)

	if err != nil {
		log.Fatalln(err)
	}

	var rules []AclRule
	for route, methods := range aclIncomingMess {
		rules = append(rules, AclRule{route, methods})
	}

	go func(ctx context.Context) error {
		lis, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				lis.Close()
				server.Stop()

				return nil
			default:
				RegisterBizServer(server, biz)

				err = server.Serve(lis)
				if err != nil {
					log.Println("Cant serve: ", err)
					return err
				}
			}
		}
	}(ctx)

	return nil
}
