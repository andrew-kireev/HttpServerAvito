package store

import (
	"HttpServerAvito/internal/model"
)

type BookingsRepository struct {
	store *Store
}

func (rep *BookingsRepository) GetallBookins(hotelId int) ([]model.Bookings, error){
	rows, err := rep.store.db.Query("SELECT * FROM bookings where hotel_id = $1", hotelId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]model.Bookings, 0)
	for rows.Next() {
		newBooking := model.Bookings{}
		err = rows.Scan(&newBooking.BookingId, &newBooking.HotelId, &newBooking.BeginData, &newBooking.EndData)
		bookings = append(bookings, newBooking)
	}
	return bookings, nil
}