package store

import (
	"HttpServerAvito/internal/model"
	"strings"
)

type BookingsRepository struct {
	store *Store
}

func (rep *BookingsRepository) AddBooking(booking *model.Bookings) (*model.Bookings, error) {
	if err := rep.store.db.QueryRow(
		"INSERT INTO bookings (hotel_id, begin_data, end_date) VALUES ($1, $2, $3) RETURNING booking_id",
		booking.HotelId, booking.BeginData, booking.EndData,
	).Scan(&booking.BookingId); err != nil {
	}

	return booking, nil
}

func (rep *BookingsRepository) GetAllBookings(hotelId int) ([]model.Bookings, error) {
	rows, err := rep.store.db.Query("SELECT * FROM bookings where hotel_id = $1 ORDER BY begin_data", hotelId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]model.Bookings, 0)
	for rows.Next() {
		newBooking := model.Bookings{}
		err = rows.Scan(&newBooking.BookingId, &newBooking.HotelId, &newBooking.BeginData, &newBooking.EndData)
		newBooking.BeginData = strings.Split(newBooking.BeginData, "T")[0]
		newBooking.EndData = strings.Split(newBooking.EndData, "T")[0]
		bookings = append(bookings, newBooking)
	}
	return bookings, nil
}

func (rep *BookingsRepository) DeleteBooking(bookingId int) error {
	_, err := rep.store.db.Exec("DELETE from bookings where booking_id = $1", bookingId)
	if err != nil {
		return err
	}

	return nil
}

func (rep *BookingsRepository) DeleteBookingsByHotelId(hotelId int) error {
	_, err := rep.store.db.Exec("DELETE from bookings where hotel_id = $1", hotelId)
	if err != nil {
		return err
	}
	return nil
}
