/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a ServerImpl for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/zencore/helloworld/proto/helloworld"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The ServerImpl port")
)

// ServerImpl is used to implement helloworld.GreeterServer.
type ServerImpl struct {
	pb.UnimplementedChatServiceServer
	store map[string]*pb.Message // store the messages.
}

func NewServer() *ServerImpl {
	result := &ServerImpl{}
	result.store = make(map[string]*pb.Message)
	return result
}

// Add a message to the store.
func (s *ServerImpl) CommitMessage(ctx context.Context, in *pb.CommitMessageRequest) (*pb.CommitMessageResponse, error) {
	log.Printf("Received commit message: %v", in.Message)
	s.store[in.Message.MessageUuid] = in.Message
	return &pb.CommitMessageResponse{}, nil
}

func (s *ServerImpl) ListMessages(ctx context.Context, in *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	// return the messages
	result := pb.ListMessagesResponse{
		Messages: []*pb.Message{},
	}
	for _, v := range s.store {
		result.Messages = append(result.Messages, v)
	}
	return &result, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &ServerImpl{})
	log.Printf("ServerImpl listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
