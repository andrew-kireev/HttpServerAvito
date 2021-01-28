package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	db *sql.DB
	hotelRep *HotelsRepository
}


func NewStore(config *Config) *Store{
	return &Store{
		config: config,
	}
}

func (store *Store) Open() error {
	db, err := sql.Open("postgres", store.config.DataBaseUrl)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	store.db = db

	return nil
}

func (store *Store) Close() {

}

func (store *Store) Hotels() *HotelsRepository  {
	if store.hotelRep != nil {
		return store.hotelRep
	}

	store.hotelRep = &HotelsRepository{
		store: store,
	}
	return store.hotelRep
}