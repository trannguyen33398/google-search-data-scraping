package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nicholasjackson/env"
	handlers "github.com/trannguyen33398/google-search-data-scraping/pkg/handler"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var connection *websocket.Conn
var connectionMutex = make(chan struct{}, 1)

type ConnectionHandler struct{}

func (ch ConnectionHandler) SendMessageToClients(message []byte) {
	connectionMutex <- struct{}{}

	err := connection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Failed to send message:", err)
		// Optionally, you can remove the connection from the map if it's no longer active

	}

	<-connectionMutex
}
func main() {

	env.Parse()

	l := log.New(os.Stdout, "data-scraping ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewUpload(l)
	
	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// getRouter := sm.Methods(http.MethodGet).Subrouter()
	// getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/upload", ph.UploadFile)
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/history-upload", ph.GetHistoryUpload)
	// Mount the socket.io server to the CORS-enabled router
	sm.HandleFunc("/ws", socketHandler)
	connectionHandler := ConnectionHandler{}
	handlers.SetConnectionHandler(connectionHandler)
	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	connectionMutex <- struct{}{}
	connection = conn
	<-connectionMutex

	for {
		// Read message from the client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		// Process the received message (e.g., log it)
		log.Println("Received message:", string(message))

		// Write a response back to the client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Failed to write message:", err)
			break
		}
	}

}
