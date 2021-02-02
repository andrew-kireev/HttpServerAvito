package store

import (
	"HttpServerAvito/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type HotelsRepository struct {
	store *Store
}

func (rep *HotelsRepository) AddHotel(h *model.Hotels) (*model.Hotels, error) {
	if err := rep.store.db.QueryRow(
		"INSERT INTO hotels (description, cost) VALUES ($1, $2) RETURNING id", h.Description, h.Price,
	).Scan(&h.Id); err != nil {
	}

	return h, nil
}

func (rep *HotelsRepository) GetHotelsList(sort string) ([]model.Hotels, error) {
	var rows *sql.Rows
	var err error
	if sort == "date" {
		rows, err = rep.store.db.Query("SELECT * FROM hotels order by creation_date")
	} else if sort == "-date" {
		rows, err = rep.store.db.Query("SELECT * FROM hotels order by creation_date desc")
	} else if sort == "price" {
		rows, err = rep.store.db.Query("SELECT * FROM hotels order by cost")
	} else if sort == "-price" {
		rows, err = rep.store.db.Query("SELECT * FROM hotels order by cost desc")
	} else {
		return nil, errors.New("not incorrect sorting param")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hotels := make([]model.Hotels, 0)
	for rows.Next() {
		newHotel := model.Hotels{}
		err = rows.Scan(&newHotel.Id, &newHotel.Description, &newHotel.Price, &newHotel.CreationDate)
		fmt.Println(newHotel)
		hotels = append(hotels, newHotel)
	}
	return hotels, nil
}

func (rep *HotelsRepository) DeleteHotel(id int) error {
	_, err := rep.store.db.Exec("DELETE from hotels where id= $1", id)
	if err != nil {
		return err
	}

	return nil
}
