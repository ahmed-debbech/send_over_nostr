package main

import (
	"errors"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Relay struct {
	Scheme     string
	Host       string
	Path       string
	Connection *websocket.Conn
	DoneCh     chan struct{}
}

func (r *Relay) Connect() error {

	serverURL := url.URL{Scheme: r.Scheme, Host: r.Host, Path: r.Path}
	r.DoneCh = make(chan struct{})

	log.Printf("connecting to %s\n", serverURL.String())
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		return errors.New("ERROR: could not dial the relay server: " + err.Error())
	}
	r.Connection = conn
	log.Println("connected to", r.Host, "successfully")
	return nil
}

func (r *Relay) Subscribe() chan []byte {

	ch := make(chan []byte)
	go func() {
		defer close(r.DoneCh)
		defer r.Connection.Close()
		for {
			_, message, err := r.Connection.ReadMessage()
			if err != nil {
				log.Println("could not read from websocket:", err)
				return
			}
			ch <- message
		}
	}()
	return ch
}

func (r *Relay) Send(message []byte) {

	select {
	case <-r.DoneCh:
		log.Println("could not send data to websocket because connection is now closed")
		return
	default:
	}

	err := r.Connection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("could not send data to websocket: " + err.Error())
		return
	}
}
