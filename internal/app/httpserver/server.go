package httpserver

import (
	"HttpServerAvito/internal/model"
	"HttpServerAvito/store"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type server struct {
	Conf   *Config
	router *http.ServeMux
	store  *store.Store
	logger *logrus.Logger
}

func NewServer(config *Config) (*server, error) {
	serv := &server{
		Conf:   config,
		router: &http.ServeMux{},
		logger: logrus.New(),
	}

	err := serv.ConfigLogger()
	if err != nil {
		return nil, err
	}

	serv.ConfigRouter()
	if err = serv.ConfigStore(); err != nil {
		serv.logger.Error("error creating db")
		return nil, err
	}
	serv.logger.Info("server created")
	return serv, nil
}

func (serv *server) HandleGetAllHotels(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("getAll Handler")
	hotels, err := serv.store.Hotels().GetHotelsList()
	if err != nil {
		serv.logger.Errorf("error in get all hotels: %v", err)
	}

	response, err := json.Marshal(hotels)
	if err != nil {
		serv.logger.Errorf("error in marshaling json: %v", err)
		return
	}

	serv.logger.Info(string(response))
	w.Write(response)
}

func (serv *server) HandleDeleteHotel(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandleDeleteHotel")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		serv.logger.Errorf("error in HandleDeleteHotel: %v", err)
		w.Write([]byte("{\"result\": \"failed\"}"))
		return
	}

	err = serv.store.Bookings().DeleteBookingsByHotelId(id)
	if err != nil {
		serv.logger.Errorf("error deleting bookings by hotel id: %v", err)
		w.Write([]byte("{\"result\": \"failed\"}"))
	}

	err = serv.store.Hotels().DeleteHotel(id)
	if err != nil {
		serv.logger.Errorf("error deleting hotels: %v", err)
		w.Write([]byte("{\"result\": \"failed\"}"))
	}
	w.Write([]byte("{\"result\": \"successful\"}"))
}

func (serv *server) HandleAddHotel(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandleAddHotel")
	hotel := &model.Hotels{}

	hotel.Description = r.FormValue("description")
	hotel.Price, _ = strconv.Atoi(r.FormValue("price"))
	hotel, err := serv.store.Hotels().AddHotel(hotel)
	if err != nil {
		serv.logger.Errorf("error in adding hotel: %v", err)
		response := fmt.Sprintf("{\"hotel_id\": %v}", 0)
		w.Write([]byte(response))
		return
	}
	response := fmt.Sprintf("{\"hotel_id\": %v}", hotel.Id)
	fmt.Println(hotel)
	w.Write([]byte(response))
}

func (serv *server) HandlerGetAllBookings(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandlerGerAllBookings")
	hotelId, _ := strconv.Atoi(r.FormValue("hotel_id"))

	bookings, err := serv.store.Bookings().GetAllBookings(hotelId)
	if err != nil {
		serv.logger.Errorf("error in getting all bookings: %v", err)
	}

	response, err := json.Marshal(bookings)
	if err != nil {
		serv.logger.Errorf("error in marshaling json: %v", err)
		return
	}

	serv.logger.Info(string(response))
	w.Write(response)
}

func (serv *server) HandlerDeleteBooking(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandlerDeleteBooking")
	bookingId, _ := strconv.Atoi(r.FormValue("booking_id"))
	err := serv.store.Bookings().DeleteBooking(bookingId)

	if err != nil {
		serv.logger.Errorf("error deleting bookings: %v", err)
		w.Write([]byte("{\"result\": \"failed\"}"))
	}
	w.Write([]byte("{\"result\": \"successful\"}"))
}

func (serv *server) HandleBookingAdd(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandleBookingAdd")
	booking := &model.Bookings{}

	booking.HotelId, _ = strconv.Atoi(r.FormValue("hotel_id"))
	booking.BeginData = r.FormValue("date_start")
	booking.EndData = r.FormValue("date_end")

	booking, err := serv.store.Bookings().AddBooking(booking)
	if err != nil {
		serv.logger.Errorf("error in adding booking: %v", err)
		response := fmt.Sprintf("{\"booking_id\": %v}", 0)
		w.Write([]byte(response))
		return
	}
	response := fmt.Sprintf("{\"booking_id\": %v}", booking.BookingId)
	fmt.Println(booking)
	w.Write([]byte(response))
}

func (serv *server) ConfigRouter() {
	serv.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serv.logger.Info("default handler")
		w.Write([]byte("default"))
	})

	serv.router.HandleFunc("/hotels/list", serv.HandleGetAllHotels)
	serv.router.HandleFunc("/hotels/delete/", serv.HandleDeleteHotel)
	serv.router.HandleFunc("/hotels/create", serv.HandleAddHotel)
	serv.router.HandleFunc("/bookings/list/", serv.HandlerGetAllBookings)
	serv.router.HandleFunc("/bookings/delete", serv.HandlerDeleteBooking)
	serv.router.HandleFunc("/bookings/create", serv.HandleBookingAdd)
}

func (serv *server) ConfigLogger() error {
	level, err := logrus.ParseLevel(serv.Conf.LogLevel)
	if err != nil {
		return err
	}

	serv.logger.SetLevel(level)
	return nil
}

func Start(config *Config) error {
	serv, err := NewServer(config)
	if err != nil {
		return err
	}

	return http.ListenAndServe(config.BindAddr, serv)
}

func (serv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serv.router.ServeHTTP(w, r)
}

func (serv *server) ConfigStore() error {
	serv.logger.Info(serv.Conf.StoreConfig.DataBaseUrl)
	st := store.NewStore(serv.Conf.StoreConfig)

	if err := st.Open(); err != nil {
		return err
	}
	serv.store = st

	return nil
}
