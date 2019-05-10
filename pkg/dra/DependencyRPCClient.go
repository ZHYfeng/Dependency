package dra

import (
	"context"
	"github.com/google/syzkaller/pkg/log"
	"google.golang.org/grpc"
	"time"
)

type DRPCClient struct {
	c DependencyRPCClient
}

func (d *DRPCClient) RunDependencyRPCClient(address *string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dependency gRPC did not connect: %v", err)
	}
	log.Logf(0, "Dependency gRPC did not connect: %v", err)
	d.c = NewDependencyRPCClient(conn)
}

func (d *DRPCClient) GetDependencyInput(name string) *NewDependencyInput {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dInput, err := d.c.GetDependencyInput(ctx, &Empty{Name: name})
	if err != nil {
		log.Fatalf("could not GetDependencyInput: %v", err)
	}
	return dInput
}

func (d *DRPCClient) SendInput(input *Input) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := d.c.SendInput(ctx, input)
	if err != nil {
		log.Fatalf("could not GetDependencyInput: %v", err)
	}
}
