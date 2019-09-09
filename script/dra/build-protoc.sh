#!/bin/bash

protoc -I=. --python_out=. ./DependencyRPC.proto