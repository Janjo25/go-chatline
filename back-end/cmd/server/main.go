package main

import (
	"fmt"
	"github.com/Janjo25/go-chatline/internal/hub"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type connectionHandler struct {
	upgrader websocket.Upgrader
	hub      *hub.Hub
}

func (handler *connectionHandler) listenForMessages(connection *websocket.Conn) {
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Error al leer mensaje: %v", err)

			break
		}

		handler.hub.Broadcast <- message // Enviar el mensaje al hub para retransmitirlo.
	}
}

func (handler *connectionHandler) handleConnections(writer http.ResponseWriter, request *http.Request) {
	connection, err := handler.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Printf("Error al actualizar la conexión: %v", err)

		return
	}
	defer func(connection *websocket.Conn) {
		err = connection.Close()
		if err != nil {
			log.Printf("Error al cerrar la conexión: %v", err)
		}
	}(connection)

	handler.hub.Register <- connection
	log.Println("Se ha establecido una nueva conexión.")

	handler.listenForMessages(connection)
}

func main() {
	h := hub.NewHub()
	go h.Run()

	handler := &connectionHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(request *http.Request) bool {
				return true
			},
		},
		hub: h,
	}

	http.HandleFunc("/ws", handler.handleConnections)

	fmt.Println("Servidor iniciado en el puerto 8080. 🚀")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
