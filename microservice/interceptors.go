package main

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	bizServer, ok := info.Server.(Biz)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	consumer, err := getConsumerFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	err = checkConsumerAccess(consumer, info.FullMethod, bizServer.rules)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func authStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	admServer, ok := srv.(Admin)
	if !ok {
		return status.Errorf(codes.Internal, "Server error")
	}

	consumer, err := getConsumerFromRequest(stream.Context())
	if err != nil {
		return err
	}

	err = checkConsumerAccess(consumer, info.FullMethod, admServer.rules)
	if err != nil {
		return err
	}

	return nil
}

func getConsumerFromRequest(ctx context.Context) ([]string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	consumer, ok := md["consumer"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Field not exist")
	}

	return consumer, nil
}

func checkConsumerAccess(consumer []string, method string, rules []AclRule) error {
	hasAccess, err := hasAccess(strings.Join(consumer, ","), method, rules)
	if err != nil {
		return err
	}

	if !hasAccess {
		return status.Errorf(codes.Unauthenticated, "Access denied")
	}

	return nil
}
