package chat

import (
	context "context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leogoesger/goservices/foundation/database"
	"github.com/pkg/errors"
)

// Chat implements ...
type Chat struct {
	DB  *sqlx.DB
	log *log.Logger
	UnimplementedChatServer
}

// New creates chat client
func New(db *sqlx.DB, log *log.Logger) *Chat {
	return &Chat{DB: db, log: log}
}

// ReadMsg ...
func (chat *Chat) ReadMsg(ctx context.Context, in *ReadRequest) (*ReadResponse, error) {

	messages := []*Message{}
	q := `SELECT * FROM messages;`

	chat.log.Printf("%s: %s", "messages.select", database.Log(q))
	if err := chat.DB.Select(&messages, q); err != nil {
		return nil, errors.Wrap(err, "select messages")
	}

	return &ReadResponse{Messages: messages}, nil
}

// WriteMsg ...
func (chat *Chat) WriteMsg(ctx context.Context, in *WriteRequest) (*Message, error) {
	now := time.Now()
	q := `INSERT INTO messages 
		(message_id, message, to_user, from_user, date_created, date_updated) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	message := Message{
		MessageId:   uuid.New().String(),
		Message:     in.Message,
		ToUser:      in.ToUser,
		FromUser:    in.FromUser,
		DateCreated: now.Format(time.RFC3339),
		DateUpdated: now.Format(time.RFC3339),
	}
	fmt.Println(in)
	chat.log.Printf("%s: %s", "messages.create", database.Log(q, message.MessageId, message.Message, message.ToUser, message.FromUser, message.DateCreated, message.DateUpdated))

	if _, err := chat.DB.ExecContext(ctx, q, message.MessageId, message.Message, message.ToUser, message.FromUser, message.DateCreated, message.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting user")
	}

	return &message, nil
}

// StreamLstMsg ...
func (chat *Chat) StreamLstMsg(_ *ReadRequest, stream Chat_StreamLstMsgServer) error {
	lastMsgID := ""

	for {
		time.Sleep(time.Millisecond * 300)
		messages := []*Message{}
		q := `SELECT * FROM messages ORDER BY date_created DESC LIMIT 1;`
		if err := chat.DB.Select(&messages, q); err != nil {
			return errors.Wrap(err, "Select last message")
		}

		if lastMsgID != messages[0].MessageId {
			if err := stream.Send(messages[0]); err != nil {
				return errors.Wrap(err, "Stream last message")
			}
			lastMsgID = messages[0].MessageId
		}
	}

}
