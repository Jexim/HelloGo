package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Jexim/HelloGo/internal/svc/hello"
)

// Migrate runs schema migrations
func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20250524_create_tables",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&hello.Hello{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("hello")
			},
		},
	})
	return m.Migrate()
}
