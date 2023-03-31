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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"proctice/configs/etcd"
	pb "proctice/proto/pb"

	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)
var serviceConfig = `{
    "loadBalancingPolicy": "round_robin"
}`

func main() {
	flag.Parse()
	// 后续改为配置文件
	etcdAddrs := []string{"127.0.0.1:2379"}
	r := etcd.NewResolver(etcdAddrs, zap.NewNop())
	resolver.Register(r)

	conn, err := grpc.Dial("etcd:///server",
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Fatalf("failed to dial %v", err)
	}
	defer conn.Close()
	fmt.Println(conn.GetState().String())
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", res.GetMessage())

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	again, err := c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "qiang"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(again.GetMessage())
}
