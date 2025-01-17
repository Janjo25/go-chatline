package hub

import (
	"github.com/gorilla/websocket"
	"log"
)

type Hub struct {
	Clients    map[*websocket.Conn]bool // Mapa de clientes conectados.
	Broadcast  chan []byte              // Canal para enviar mensajes a todos los clientes.
	Register   chan *websocket.Conn     // Canal para registrar nuevas conexiones.
	Unregister chan *websocket.Conn     // Canal para desconectar clientes.
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

// Run inicia el bucle principal del hub. Se encarga de manejar las conexiones y los mensajes que se envían.
func (hub *Hub) Run() {
	for {
		select {
		case connection := <-hub.Register:
			hub.Clients[connection] = true
			log.Println("Nuevo cliente registrado")
		case connection := <-hub.Unregister:
			if _, ok := hub.Clients[connection]; ok {
				delete(hub.Clients, connection)

				err := connection.Close()
				if err != nil {
					log.Printf("Error al cerrar la conexión al intentar desconectar al cliente: %v", err)
				}

				log.Println("Cliente desconectado")
			}
		case message := <-hub.Broadcast:
			for connection := range hub.Clients {
				err := connection.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error al enviar mensaje: %v", err)

					err = connection.Close()
					if err != nil {
						log.Printf("Error al cerrar la conexión luego de fallar al enviar mensaje: %v", err)
					}

					delete(hub.Clients, connection)
				}
			}
		}
	}
}
