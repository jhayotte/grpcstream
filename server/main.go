package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/jhayotte/grpcstream/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port      = ":50051"
	chunkSize = 200 // In a real sample, we should try to choose a chunk of a size of 64 KiB due to https://github.com/grpc/grpc.github.io/issues/371
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

//inventory of all existing book
var books []*pb.Book

func NewBook() {
	var i int64
	for {
		books = append(books, &pb.Book{
			Author:      "JK",
			Title:       fmt.Sprintf("book %d", i),
			Description: "",
		})
		i++
		if i > 3000 {
			break
		}
	}
}

func (s *server) GetAllBooksByAuthor(req *pb.GetAllBooksByAuthorRequest, srv pb.BookService_GetAllBooksByAuthorServer) (err error) {
	var start int64
	for {
		end := start + chunkSize
		if int(end) > len(books) {
			end = int64(len(books))
		}
		if int(start) >= len(books) {
			break
		}
		err = srv.Send(&pb.GetAllBooksByAuthorResponse{
			Books: books[start:end],
		})

		if err != nil {
			return errors.New(fmt.Sprintf("cannot send message: %v", err))
		}
		time.Sleep(10 * time.Millisecond)
		start += chunkSize
		if end > int64(len(books)) {
			// we reached the end of our message.
			// we close the channel.
			return nil
		}
	}
}

func main() {
	// initialize our book inventory
	NewBook()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBookServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
