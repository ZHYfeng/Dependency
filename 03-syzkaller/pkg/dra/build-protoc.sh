#!/bin/bash

protoc ./DependencyRPC.proto  --go_out=plugins=grpc:.