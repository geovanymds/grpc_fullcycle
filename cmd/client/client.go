package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/geovanymds/grpc_fullcycle/pb/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	//AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Geovany",
		Email: "geovany.mendes@estantemagica.com.br",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Geovany",
		Email: "geovany.mendes@estantemagica.com.br",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {

		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not make gRPC request: %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "Geo",
			Email: "geo@mendes.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "Geo 2",
			Email: "geo2@mendes.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "Geo 3",
			Email: "geo3@mendes.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "Geo 4",
			Email: "geo4@mendes.com",
		},
		&pb.User{
			Id:    "g5",
			Name:  "Geo 5",
			Email: "geo5@mendes.com",
		},
		&pb.User{
			Id:    "g6",
			Name:  "Geo 6",
			Email: "geo6@mendes.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "Geo",
			Email: "geo@mendes.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "Geo 2",
			Email: "geo2@mendes.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "Geo 3",
			Email: "geo3@mendes.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "Geo 4",
			Email: "geo4@mendes.com",
		},
		&pb.User{
			Id:    "g5",
			Name:  "Geo 5",
			Email: "geo5@mendes.com",
		},
		&pb.User{
			Id:    "g6",
			Name:  "Geo 6",
			Email: "geo6@mendes.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
