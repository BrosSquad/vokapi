package container

import (
	"context"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type (
	Config struct {
		BadgedDBPath string `json:"badger_db_path"`
	}
	Container struct {
		ctx      context.Context
		config   *Config
		badgerDB *badger.DB
	}
)

func New(ctx context.Context, cfg *Config) *Container {
	return &Container{
		ctx:    ctx,
		config: cfg,
	}
}

func (c *Container) GetBadgerDB() *badger.DB {
	if c.badgerDB == nil {
		options := badger.DefaultOptions(c.config.BadgedDBPath).
			WithIndexCacheSize(250 << 20).
			WithLoggingLevel(badger.ERROR)

		db, err := badger.Open(options)
		if err != nil {
			log.Fatalf("Faield to open badger db: %v", err)
		}

		go func(db *badger.DB) {
			ticker := time.NewTicker(5 * time.Minute)
			for {
				select {
				case <-ticker.C:
					err := db.RunValueLogGC(0.7)
					if err != nil {
						log.Printf("Failed to run value log GC: %v\n", err)
					}
				case <-c.ctx.Done():
					return
				}
			}
		}(db)

		c.badgerDB = db
	}

	return c.badgerDB
}

func (c *Container) Close() error {
	return nil
}
