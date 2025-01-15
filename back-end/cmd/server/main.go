package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (handler *webSocketHandler) handleConnections(writer http.ResponseWriter, request *http.Request) {
	connection, err := handler.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("Error al actualizar la conexión:", err)

		return
	}
	defer func(connection *websocket.Conn) {
		err = connection.Close()
		if err != nil {
			log.Println("Error al cerrar la conexión:", err)
		}
	}(connection)

	log.Println("Se ha establecido una nueva conexión.")
}

func main() {
	handler := &webSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(request *http.Request) bool {
				return true
			},
		},
	}

	http.HandleFunc("/ws", handler.handleConnections)

	fmt.Println("Servidor iniciado en el puerto 8080. 🚀")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
