package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"text/template"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	_ "github.com/cirius-go/portfolio-server/cmd/workers/migrate/migrations"
	"github.com/cirius-go/portfolio-server/internal/config"
	"github.com/cirius-go/portfolio-server/pkg/db"
)

var (
	//go:embed migrations
	migrations embed.FS
	cfgFile    = flag.String("cfg", ".env", "the path to the config file")

	tmpl = template.Must(template.New("goose.gorm-migration").Parse(`
package migrations

import (
  "context"
  "database/sql"

  "github.com/pressly/goose/v3"
  "gorm.io/gorm"
)

func init() {
  goose.AddMigrationNoTxContext(up{{.CamelName}}, down{{.CamelName}})
}

func up{{.CamelName}}(ctx context.Context, tx *sql.DB) error {
  gdb, err := initDB(tx)
  if err != nil {
    return err
  }

  return gdb.Transaction(func(tx *gorm.DB) error {
    // return tx.Migrator().AutoMigrate()
    return execSlice(tx)
  })
}

func down{{.CamelName}}(ctx context.Context, tx *sql.DB) error {
  gdb, err := initDB(tx)
  if err != nil {
    return err
  }

  return gdb.Transaction(func(tx *gorm.DB) error {
    // return tx.Migrator().DropTable()
    return execSlice(tx)
  })
}`))
)

func main() {
	flag.Parse()

	cfg, err := config.Load(*cfgFile)
	if err != nil {
		panic(err)
	}

	pg, err := db.NewPostgres(cfg.PGDB)
	if err != nil {
		panic(err)
	}
	defer pg.Conn.Close()

	goose.SetBaseFS(migrations)
	goose.SetDialect("postgres")
	goose.SetTableName("migrations")
	goose.SetVerbose(true)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	root := cobra.Command{
		Use:               "migrate",
		Short:             "Migrate database schema",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
	}

	root.AddCommand(&cobra.Command{
		Use:     "create",
		Short:   `Create a new migration file`,
		Example: `migrate create gorm create_customer_table`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "gorm":
				return goose.CreateWithTemplate(pg.Conn, "cmd/workers/migrate/migrations", tmpl, args[1], "go")
			default:
				return errors.New("unimplemented")
			}
		},
	})

	root.AddCommand(&cobra.Command{
		Use:     "up",
		Short:   "upgrade to the latest version",
		Example: "migrate up",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return goose.RunWithOptionsContext(ctx, "up", pg.Conn, "migrations", args, goose.WithAllowMissing())
		},
	})

	root.AddCommand(&cobra.Command{
		Use:     "down",
		Short:   "downgrade to the previous version",
		Example: `migrate down`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return goose.RunWithOptionsContext(ctx, "down", pg.Conn, "migrations", args, goose.WithAllowMissing())
		},
	})

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
