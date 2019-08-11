package main

import (
	"context"
	"log"
)

type Biz struct {
	rules []AclRule
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
	return &Nothing{Dummy: true}, nil
}
