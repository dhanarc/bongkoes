package shared

import (
	"github.com/djk-lgtm/atlassianoto/config"
	"github.com/djk-lgtm/atlassianoto/internal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(c *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(c.Atlassianoto.DBLocation), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(internal.DAOs...)
	if err != nil {
		panic("failed to migrate database")
	}
	return db
}
