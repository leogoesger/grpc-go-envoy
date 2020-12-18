// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1.1,
		Description: "Create table messages",
		Script: `
		CREATE TABLE messages (
			message_id       UUID UNIQUE,
			to_user          TEXT,
			from_user        TEXT,
			message          TEXT,
			date_created  	 TIMESTAMP,
			date_updated  	 TIMESTAMP
		);`,
	},
}
