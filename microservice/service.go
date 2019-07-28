package main

import (
	"context"
	fmt "fmt"
	"log"
	"net"

	grpc "google.golang.org/grpc"
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
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Println("cant listen port: ", err)
		return err
	}

	server := grpc.NewServer()
	biz := Biz{}

	RegisterBizServer(server, biz)
	fmt.Println("starting server at :8082")
	err = server.Serve(lis)
	if err != nil {
		log.Println("Cant serve: ", err)
		return err
	}

	return nil
}
