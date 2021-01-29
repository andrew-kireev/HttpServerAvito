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

func (serv *server) HandleDeleteHotel(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandleDeleteHotel")
	fmt.Println(r.URL)
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		serv.logger.Errorf("error in HandleDeleteHotel: %v", err)
		w.Write([]byte("Произошла ошибка при удалении"))
		return
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
		response:= fmt.Sprintf("{\"hotel_id\": %v}", 0)
		w.Write([]byte(response))
		return
	}
	response:= fmt.Sprintf("{\"hotel_id\": %v}", hotel.Id)
	fmt.Println(hotel)
	w.Write([]byte(response))
}

func (serv *server) HandlerGerAllBookings(w http.ResponseWriter, r *http.Request) {
	serv.logger.Info("HandlerGerAllBookings")
	hotelId, _ := strconv.Atoi(r.FormValue("hotel_id"))

	bookings, err := serv.store.Bookings().GetallBookins(hotelId)
	if err != nil {
		serv.logger.Errorf("error in getting all bookings: %v", err)
	}

	fmt.Println(bookings)
}

func (serv *server) ConfigRouter() {
	serv.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serv.logger.Info("default handler")
		w.Write([]byte("default"))
	})

	serv.router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		serv.logger.Info("hello handler")
		w.Write([]byte("Hello"))
	})

	serv.router.HandleFunc("/delete/", serv.HandleDeleteHotel)
	serv.router.HandleFunc("/addHotel/", serv.HandleAddHotel)
	serv.router.HandleFunc("/bookings/list/", serv.HandlerGerAllBookings)

	serv.router.HandleFunc("/getAll", func(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("\nbye"))
	})
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
	st := store.NewStore(serv.Conf.StoreConfig)

	if err := st.Open(); err != nil {
		return err
	}
	serv.store = st

	return nil
}
