package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin:     func(r *http.Request) bool { return true },
}

type App struct {
	hub  *Hub
	conf *ServerConf
}

// Run the application server
func (a *App) Run() error {
	go a.hub.run()
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.RootHandler)
	mux.HandleFunc("/ws", a.WSHandler)
	if err := http.ListenAndServe(a.conf.Host, mux); err != nil {
		return fmt.Errorf(`server failed to serve: %s`, err)
	}
	return nil
}

func NewApp(conf *ServerConf) *App {
	return &App{
		hub:  NewHub(),
		conf: conf,
	}
}

func (a App) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func (a App) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade: %s", err)
		http.Error(w, "failed to upgrade", 500)
		return
	}

	client := newClient(a.hub, conn)
	a.hub.addClient(client)
	client.run()
}
