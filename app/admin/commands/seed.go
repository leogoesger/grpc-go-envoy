package commands

import (
	"fmt"

	"github.com/leogoesger/goservices/foundation/database"
	"github.com/leogoesger/grpc-go/business/schema"
	"github.com/pkg/errors"
)

// Seed loads test data into the database.
func Seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect database")
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return errors.Wrap(err, "seed database")
	}

	fmt.Println("seed data complete")
	return nil
}
