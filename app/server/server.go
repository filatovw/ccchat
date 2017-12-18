package server

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/filatovw/ccchat/app/server/model"
	"github.com/filatovw/ccchat/internal/protocol"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type App struct {
	hub  *Hub
	conf *Conf
}

// Run the application server
func (a *App) Run() error {
	dbconn, err := model.InitDB(a.conf.DBHost, a.conf.DBUser, a.conf.DBPass, a.conf.DBName)
	if err != nil {
		return errors.Wrap(err, "failed to establish new connection to database")
	}
	defer dbconn.Close()
	a.hub.setDB(dbconn)
	go a.hub.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.RootHandler)
	mux.HandleFunc("/static/", a.StaticHandler)
	mux.HandleFunc("/ws", a.WSHandler)
	if err := http.ListenAndServe(a.conf.Host, mux); err != nil {
		return errors.Wrap(err, `server failed to start`)
	}
	return nil
}

func NewApp(conf *Conf) *App {
	return &App{
		hub:  NewHub(),
		conf: conf,
	}
}

func (a App) StaticHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("StaticHandler: %s", r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, "/static") {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := Asset(strings.Join([]string{"app/server", r.URL.Path}, ""))
	if err != nil {
		log.Printf(`failed to get file from bindata: %s`, err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	splitted := strings.Split(r.URL.Path, ".")
	var contentType string
	if len(splitted) > 1 {
		switch splitted[len(splitted)-1]{
		case "js":
			contentType = "application/json"
		case "css":
			contentType = "text/css"
		default:
			contentType = "text/plain"
		}
	} else {
		contentType = "text/plain"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (a App) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	index, err := Asset("app/server/static/tpl/index.html")
	if err != nil {
		log.Printf(`failed to load template from bindata: %s`, err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	t := template.New("index")
	t, err = t.Parse(string(index))
	if err != nil {
		log.Printf(`failed to parse template from bindata: %s`, err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	// collect messages for a day
	yesterday := time.Now().AddDate(0, 0, -1)
	messages, err := model.MessagesSinceDate(a.hub.db, yesterday)
	if err != nil {
		log.Printf(`failed to get messages: %s`, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	messagesT := []string{}
	for _, m := range messages {
		msg, err := protocol.ParseMessage([]byte(m.MessageBody))
		if err != nil {
			log.Printf("%s", err)
		}
		messagesT = append(messagesT, string(msg.MarshalServer(m.UserName)))
	}

	err = t.Execute(w, messagesT)
	if err != nil {
		log.Printf(`failed to render template: %s`, err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
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
