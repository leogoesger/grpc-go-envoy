package chat

import (
	context "context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Chat implements ...
type Chat struct {
	DB *sqlx.DB
	UnimplementedChatServer
}

// New creates chat client
func New(db *sqlx.DB) *Chat {
	return &Chat{DB: db}
}

// ReadMsg ...
func (chat *Chat) ReadMsg(ctx context.Context, in *ReadRequest) (*ReadResponse, error) {

	messages := []*Message{}
	q := `SELECT * FROM messages;`

	if err := chat.DB.Select(&messages, q); err != nil {
		return nil, errors.Wrap(err, "select messages")
	}

	return &ReadResponse{Messages: messages}, nil
}

// WriteMsg ...
func (chat *Chat) WriteMsg(ctx context.Context, in *WriteRequest) (*ReadResponse, error) {
	messages := []*Message{
		{Message: "hello", ToUser: "Leo", FromUser: "Noelle"},
		{Message: "hello2", ToUser: "Leo2", FromUser: "Noelle2"},
		{Message: "hello3", ToUser: "Leo3", FromUser: "Noelle3"},
	}
	return &ReadResponse{Messages: messages}, nil
}
