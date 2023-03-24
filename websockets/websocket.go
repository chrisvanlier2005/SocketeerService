package websockets

import (
	"JobHiraMicroservice/models"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"log"
	"math/rand"
	"strconv"
	"sync"
)

type client struct {
	id        string
	clientKey string
	isClosing bool
	mu        sync.Mutex
}

type IncomingConnection struct {
	Key        string
	Connection *websocket.Conn
}

type BroadcastMessageType struct {
	Message   string
	ServerKey string
	ClientKey string
}

type Channel struct {
	ChannelName string
}

var Clients = make(map[*websocket.Conn]*client)
var Register = make(chan *websocket.Conn)

// array of key and connection
var RegisterWithKey = make(chan IncomingConnection)
var Broadcast = make(chan string)
var Unregister = make(chan *websocket.Conn)

func BroadcastMessage(broadCastMessage BroadcastMessageType, channel Channel) error {
	application := models.Application{}
	if err := models.DB.Where("server_key = ?", broadCastMessage.ServerKey).First(&application).Error; err != nil {
		fmt.Println(err)
		return err
	}
	for connection, c := range Clients {
		go func(connection *websocket.Conn, c *client) { // send to each client in parallel so we don't block on a slow client
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.isClosing {
				return
			}
			if c.clientKey != application.ClientKey {
				return
			}
			fmt.Println("sending message to client that matches: ", c.id)
			if err := connection.WriteMessage(websocket.TextMessage, []byte(broadCastMessage.Message)); err != nil {
				c.isClosing = true
				log.Println("write error:", err)

				connection.WriteMessage(websocket.CloseMessage, []byte{})
				connection.Close()
				Unregister <- connection
			}
		}(connection, c)
	}
	return nil
}

func RunHub() {
	for {
		select {
		// register a new connection.
		case incomingConnection := <-RegisterWithKey:
			clientKey := incomingConnection.Key
			connection := incomingConnection.Connection

			clientConnection := &client{
				id:        strconv.Itoa(rand.Int()),
				clientKey: clientKey,
			}

			Clients[connection] = clientConnection
			if err := connection.WriteMessage(websocket.TextMessage, []byte("Connected succesfully as "+clientConnection.id)); err != nil {
				log.Println("Error writing", err)
				connection.WriteMessage(websocket.CloseMessage, []byte{})
				connection.Close()
				Unregister <- connection
			}
			fmt.Println("client registered with key: ", clientKey)

		// Upon message received
		case message := <-Broadcast:
			log.Println("message received:", message)
			// Send the message to all clients

		// Unregistering the user.
		case connection := <-Unregister:
			// Remove the client from the hub
			delete(Clients, connection)

			log.Println("connection unregistered")
		}
	}
}
