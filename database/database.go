package database

import (
	"fmt"
	"net/url"

	"github.com/omhen/swissblock-trade-executor/v2/errors"
	"github.com/omhen/swissblock-trade-executor/v2/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ContextKeyDB = "db"
)

func NewDBConnection(dbURL string) (*gorm.DB, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, errors.Annotate(err, "Malformed Database URL")
	}

	var engine gorm.Dialector
	var dbConfig gorm.Config
	switch u.Scheme {
	case "sqlite":
		engine = sqlite.Open(u.Path)
		dbConfig = gorm.Config{}
	default:
		log.WithFields(log.Fields{
			"databaseURL": dbURL,
		}).Error("Unsupported database url scheme")
		return nil, errors.New(fmt.Sprintf("Unsupported dabase url scheme %s", u.Scheme))
	}

	db, err := gorm.Open(engine, &dbConfig)
	if err == nil {
		db.AutoMigrate(&model.Order{}, &model.Trade{})
	}
	return db, err
}
