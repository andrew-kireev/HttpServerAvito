package httpserver

import (
	"fmt"
	"net/http"
)


type server struct {
	//router *mux.Router
	Address string
	Conf    *Config
	router http.ServeMux
}

func NewServer(config *Config) *server{
	serv := &server{
		Address: "localhost",
		Conf: config ,
	}

	serv.ConfigRouter()
	return serv
}

func (serv *server)ConfigRouter() {
	serv.router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("default"))
	})

	serv.router.HandleFunc("/hello", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	serv.router.HandleFunc("/bye", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bye"))
	})
}

func Start(config *Config) error {
	serv := NewServer(config)
	fmt.Println(serv.Address)


	return http.ListenAndServe(config.BindAddr, serv)
}

func (serv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serv.router.ServeHTTP(w, r)
}

