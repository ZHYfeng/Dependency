package dra

import (
	"context"
	"github.com/google/syzkaller/pkg/log"
	"google.golang.org/grpc"
	"time"
)

type DRPCClient struct {
	c DependencyRPCClient
	i []*Input
}

func (d *DRPCClient) RunDependencyRPCClient(address, name *string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dependency gRPC did not connect: %v", err)
	}
	d.c = NewDependencyRPCClient(conn)
	d.Connect(name)
}

func (d *DRPCClient) Connect(name *string) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := d.c.Connect(ctx, &Empty{Name: *name})
	if err != nil {
		log.Fatalf("Dependency gRPC could not Connect: %v", err)
	}
	return
}

func (d *DRPCClient) GetDependencyInput(name string) *NewDependencyInput {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dInput, err := d.c.GetDependencyInput(ctx, &Empty{Name: name})
	if err != nil {
		log.Fatalf("Dependency gRPC could not GetDependencyInput: %v", err)
	}
	return dInput
}

func (d *DRPCClient) SendInput(input *Input) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	d.i = append(d.i, input)
	if len(d.i) == 100 {
		for _, ii := range d.i {
			_, err := d.c.SendInput(ctx, ii)
			if err != nil {
				log.Fatalf("Dependency gRPC could not SendInput: %v", err)
			}
		}
		d.i = nil
	}
}
