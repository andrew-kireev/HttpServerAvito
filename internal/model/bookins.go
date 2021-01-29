package model

type Bookings struct {
	BookingId int    `json:"booking_id"`
	HotelId   int    `json:"hotel_id"`
	BeginData string `json:"begin_date"`
	EndData   string `json:"end_date"`
}
