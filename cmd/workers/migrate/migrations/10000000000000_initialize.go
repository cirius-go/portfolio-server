package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationNoTxContext(upInitialize, downInitialize)
}

var initTables = []any{}

func upInitialize(ctx context.Context, tx *sql.DB) error {
	gdb, err := initDB(tx)
	if err != nil {
		return err
	}

	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Migrator().AutoMigrate(initTables...); err != nil {
			return err
		}

		return nil
	})
}

func downInitialize(ctx context.Context, tx *sql.DB) error {
	gdb, err := initDB(tx)
	if err != nil {
		return err
	}

	return gdb.Transaction(func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(initTables...)
	})
}
