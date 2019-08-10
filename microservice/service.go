package main

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type AdminServerHandler interface {
	Logging(context.Context, Nothing) error
	Statistics(context.Context, Nothing) error
}

/*type Admin struct {
}

func (adm Admin) Logging(nothing *Nothing) error {
	log.Println("*Logging()*")
	return nil
}

func (adm Admin) Statistics(ctx context.Context, nothing *Nothing) error {
	log.Println("*Statistics()*")
	return nil
}*/

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)

	rules, err := CreateRulesFromIncomingMessage([]byte(ACLData))
	if err != nil {
		return err
	}

	biz := Biz{rules}
	//adm := Admin{}

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
				//RegisterAdminServer(server, adm)

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

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	bizServer := info.Server.(Biz)

	md, _ := metadata.FromIncomingContext(ctx)
	consumer, ok := md["consumer"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Field not exist")
	}

	hasAccess, err := hasAccess(strings.Join(consumer, ","), info.FullMethod, bizServer.rules)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, status.Errorf(codes.Unauthenticated, "Access denied")
	}

	return handler(ctx, req)
}
