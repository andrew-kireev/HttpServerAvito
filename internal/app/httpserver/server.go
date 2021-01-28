package httpserver

import (
	"HttpServerAvito/store"
	"fmt"
	"net/http"
)


type server struct {
	//router *mux.Router
	Conf    *Config
	router *http.ServeMux
	store *store.Store
}

func NewServer(config *Config) (*server, error) {
	serv := &server{
		Conf: config,
		router: &http.ServeMux{},
	}

	serv.ConfigRouter()
	if err := serv.ConfigStore(); err != nil {
		fmt.Println("error opening db")
		return nil, err
	}
	return serv, nil
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

