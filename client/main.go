package main

import (
	"log"
	"os"

	pb "github.com/jhayotte/grpcstream/book"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "JK"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBookServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.GetAllBooksByAuthor(context.Background(), &pb.GetAllBooksByAuthorRequest{Author: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		result, err := r.Recv()
		if err != nil {
			log.Fatalf("could not receive: %v", err)
		}

		for _, item := range result.Books {
			log.Printf("book: %s  %s  %s", item.Author, item.Title, item.Description)
		}

		if err = r.CloseSend(); err != nil {
			break
		}
	}
}
