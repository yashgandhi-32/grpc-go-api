#!/bin/bash

protoc blogproto/blog.proto  --go_out=plugins=grpc:.