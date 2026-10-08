package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "chat/pkg/note_v1"
)

const (
	address = "localhost:50301"
	noteId  = 142
)

func main() {

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to server :", err)
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: noteId})
	if err != nil {
		log.Fatalf("failed to get note by id:", noteId)
	}
	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetNote()))

}
