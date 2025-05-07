package migrations

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(conn *sql.DB) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		AllowGlobalUpdate:                        false,
		IgnoreRelationshipsWhenMigrating:         true,
	})
}

func execSlice(tx *gorm.DB, changes ...string) error {
	for i := range changes {
		if err := tx.Exec(changes[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
