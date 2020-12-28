package helloworld

import (
	context "context"
	"log"

	"github.com/jmoiron/sqlx"
	helloworldv1 "github.com/leogoesger/grpc-go/proto/gen/grpc-go/helloworld/v1"
)

// HelloWorld implements ...
type HelloWorld struct {
	DB  *sqlx.DB
	log *log.Logger
	helloworldv1.UnimplementedGreeterServer
}

// New creates chat client
func New(db *sqlx.DB, log *log.Logger) *HelloWorld {
	return &HelloWorld{DB: db, log: log}
}

// SayHello implements helloworld.GreeterServer
func (s *HelloWorld) SayHello(ctx context.Context, in *helloworldv1.HelloRequest) (*helloworldv1.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworldv1.HelloReply{Message: "Hello " + in.GetName()}, nil
}
