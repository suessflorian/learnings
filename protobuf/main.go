package main

import (
	"addressbook/protobuf"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var addressbook = &protobuf.AddressBook{
	People: []*protobuf.Person{
		{
			Name:  "Florian",
			Id:    0,
			Email: "floriansuess96@icloud.com",
			Phones: []*protobuf.Person_PhoneNumber{{
				Number: "02108018366",
				Type:   protobuf.Person_MOBILE,
			}},
			LastUpdated: new(timestamppb.Timestamp),
		},
	},
}

const SERVER_URL = "localhost:8080"

func main() {
	listening, err := net.Listen("tcp", SERVER_URL)
	if err != nil {
		panic(fmt.Errorf("failed to listen on port 8080: %w", err))
	}
	grpcServer := grpc.NewServer()
	protobuf.RegisterAddressBookServiceServer(grpcServer, new(server))

	go func() {
		log.Println("grpc server listening on port 8080...")

		if err := grpcServer.Serve(listening); err != nil {
			panic(fmt.Errorf("server unexpectedly shut down: %w", err))
		}
	}()

	log.Println("spinning up client...")
	conn, err := grpc.Dial(SERVER_URL, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Errorf("failed to establish client grpc connection to %s: %w", SERVER_URL, err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client := protobuf.NewAddressBookServiceClient(conn)
	var protoReqCount int
	for {
		_, err := client.GetPerson(ctx, &protobuf.GetPersonRequest{
			Id: 0,
		})
		if err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.Canceled) {
				break
			}
			panic(fmt.Errorf("failed to get person over grpc: %w", err))
		}
		protoReqCount += 1
	}

	log.Printf("shutting down grpc server...")
	grpcServer.Stop()
	log.Printf("shutting down grpc client...")
	if err := conn.Close(); err != nil {
		panic(fmt.Errorf("failed to close grpc dial connection: %w", err))
	}

	log.Printf("finished grpc excercise with %d requests\n", protoReqCount)

	http.HandleFunc("/person", func(rw http.ResponseWriter, r *http.Request) {
		reqId := r.URL.Query().Get("id")
		id, err := strconv.Atoi(reqId)
		if err != nil {
			panic(fmt.Errorf("failed to extract person id out of url: %w", err))
		}
		resp, err := getPerson(r.Context(), int32(id))
		if err != nil {
			panic(fmt.Errorf("failed to 'getPerson': %w", err))
		}

		marshalledResp, err := json.Marshal(resp.GetPerson())
		if err != nil {
			panic(fmt.Errorf("failed to marshal person response: %w", err))
		}

		if _, err := rw.Write(marshalledResp); err != nil {
			panic(fmt.Errorf("write marshaled person to response writer: %w", err))
		}
	})

	go func() {
		err := http.ListenAndServe(SERVER_URL, nil)
		if err != http.ErrServerClosed {
			panic(fmt.Errorf("unexpected http server shutdown: %w", err))
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var (
		httpReqCount     int
		httpFailureCount int
	)
	for {
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			break
		}
		_, err := http.Get("http://" + SERVER_URL + "/person?id=0")
		if err != nil {
			httpFailureCount += 1
		}

		httpReqCount += 1
	}

	log.Printf("finished http excercise with %d requests with %d failues", httpReqCount, httpFailureCount)
}

type server struct {
	protobuf.UnimplementedAddressBookServiceServer
}

func (s *server) GetPerson(ctx context.Context, req *protobuf.GetPersonRequest) (*protobuf.GetPersonResponse, error) {
	return getPerson(ctx, req.GetId())
}

func getPerson(ctx context.Context, id int32) (*protobuf.GetPersonResponse, error) {
	for _, person := range addressbook.People {
		if person.Id != id {
			continue
		}
		return &protobuf.GetPersonResponse{
			Person: person,
		}, nil
	}

	return nil, fmt.Errorf("failed to find person with id: %d", id)
}
