package store

import (
	"HttpServerAvito/internal/model"
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

func (rep *HotelsRepository) GetHotelsList() ([]model.Hotels, error) {
	rows, err := rep.store.db.Query("SELECT * FROM hotels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hotels := make([]model.Hotels, 0)
	for rows.Next() {
		newHotel := model.Hotels{}
		err = rows.Scan(&newHotel.Id, &newHotel.Description, &newHotel.Price)
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
