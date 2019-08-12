package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)

	rules, err := CreateRulesFromIncomingMessage([]byte(ACLData))
	if err != nil {
		return err
	}

	biz := Biz{rules}
	adm := Admin{rules: rules}

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
				RegisterAdminServer(server, adm)

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
