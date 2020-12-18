package helloworld

import (
	context "context"
	"log"

	"github.com/jmoiron/sqlx"
)

// HelloWorld implements ...
type HelloWorld struct {
	DB *sqlx.DB
	UnimplementedGreeterServer
}

// New creates chat client
func New(db *sqlx.DB) *HelloWorld {
	return &HelloWorld{DB: db}
}

// SayHello implements helloworld.GreeterServer
func (s *HelloWorld) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &HelloReply{Message: "Hello " + in.GetName()}, nil
}
