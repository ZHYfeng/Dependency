package dra

import (
	"context"
	"github.com/google/syzkaller/pkg/log"
	"google.golang.org/grpc"
	"time"
)

type DRPCClient struct {
	c DependencyRPCClient
	I []*Input
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

	d.I = append(d.I, cloneInput(input))
	if len(d.I) == 10 {
		for _, ii := range d.I {
			_, err := d.c.SendInput(ctx, ii)
			if err != nil {
				log.Fatalf("Dependency gRPC could not SendInput: %v", err)
			}
			log.Logf(3, "Dependency gRPC SendInput sig : v%", ii.Sig)
			for _, cc := range ii.Call {
				log.Logf(3, "Dependency gRPC SendInput address : v%", cc.Address)
			}
		}
	}

	d.I = nil
	n, _ := d.c.GetNewInput(ctx, &Empty{})
	log.Logf(3, "Dependency gRPC GetNewInput : v%", len(n.Input))
	for _, aa := range n.Input {
		log.Logf(3, "Dependency gRPC GetNewInput sig : v%", aa.Sig)
		for _, cc := range aa.Call {
			log.Logf(3, "Dependency gRPC GetNewInput address : v%", cc.Address)
		}
	}
}
