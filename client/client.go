package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"

	api "gRPC-project/api"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Не могу подключиться: %v", err)
	}
	defer conn.Close()

	client := api.NewKeyValueServiceClient(conn)

	flag := true
	for flag {
		var cmd string
		fmt.Fscan(os.Stdin, &cmd)
		if err != nil {
			return
		}

		switch cmd {
		case "put":
			var id int32
			var val string

			fmt.Fscan(os.Stdin, &id, &val)

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

			r, err := client.Put(ctx, &api.PutKeyValueRequest{Id: id, Val: val})
			if err != nil {
				fmt.Printf("error %v\n", err)
			} else {
				fmt.Printf("{%d, %v}\n", r.Id, r.Val)
			}

		case "get":
			var id int32

			fmt.Fscan(os.Stdin, &id)

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

			r, err := client.FindById(ctx, &api.GetKeyValueRequest{Id: id})
			if err != nil {
				fmt.Printf("error %v\n", err)
			} else {
				fmt.Printf("{%d, %v}\n", r.Id, r.Val)
			}
		case "delete":
			var id int32

			fmt.Fscan(os.Stdin, &id)

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

			r, err := client.Delete(ctx, &api.DeleteKeyValue{Id: id})
			if err != nil {
				fmt.Printf("error %v\n", err)
			} else {
				fmt.Printf("{%d, %v}\n", r.Id, r.Val)
			}
		case "getMany":
			var len int32
			var page int32

			fmt.Fscan(os.Stdin, &len, &page)

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

			r, err := client.ManyKeyValues(ctx, &api.PagingRequest{PageLength: len, PageNumber: page})
			if err != nil {
				fmt.Printf("error %v", err)
			} else {
				for _, kv := range r.KeyValues {
					fmt.Printf("{%d, %v}\n", kv.Id, kv.Val)
				}
			}
		default:
			flag = false
		}
	}
}
