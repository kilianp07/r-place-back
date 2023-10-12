// routes/routes.go

package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kilianp07/r-place-back/db"
	"github.com/kilianp07/r-place-back/utils"
)

const (
	Width  = 1000
	Height = 1000
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	register        = make(chan *Client)
	broadcast       = make(chan utils.Pixel)
	unregister      = make(chan *Client)
	clients         = make(map[*Client]bool)
	colorPlacements [Width][Height]utils.Pixel
	mu              sync.Mutex
)

type Route struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

type Client struct {
	Conn *websocket.Conn
	Send chan utils.Pixel
}

func (c *Client) WritePump() {
	for placement := range c.Send {
		err := c.Conn.WriteJSON(placement)
		if err != nil {
			break
		}
	}
}

var routes = []Route{
	{Path: "/websocket", Handler: WebSocketHandler, Method: "GET"},
	{Path: "/color-placements", Handler: ColorPlacementsHandler, Method: "GET"},
}

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	for _, route := range routes {
		r.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}

	go func() {
		for {
			client := <-register // Attendez qu'un client soit inscrit
			clients[client] = true
		}
	}()
	go BroadcastColorPlacements()

	return r
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")

	// Créez une structure de données pour représenter chaque client
	client := &Client{
		Conn: conn,
		Send: make(chan utils.Pixel), // Canal pour envoyer des données au client
	}

	// Inscrivez le client à la liste des clients connectés
	register <- client
	defer func() {
		// Désinscrire le client lorsque la connexion se termine
		unregister <- client
	}()

	// Lancez une goroutine pour écouter les messages du client
	go client.WritePump()

	for {
		// Lire un message du client
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var placement utils.Pixel
		if err := json.Unmarshal(p, &placement); err != nil {
			continue
		}

		// Vérifier les coordonnées pour éviter les dépassements
		if placement.X >= 0 && placement.X < Width && placement.Y >= 0 && placement.Y < Height {
			// Mettre à jour le tableau de pixels avec le nouveau placement
			colorPlacements[placement.X][placement.Y] = placement

			// Insérer les données dans la base de données

			if err := db.InsertPixel(placement); err != nil {
				fmt.Println(err)
				continue
			}

			// Envoyer le placement à tous les clients connectés
			broadcast <- placement
		}
	}
}

func ColorPlacementsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	colorPlacementsJSON, err := json.Marshal(colorPlacements)
	mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(colorPlacementsJSON)
}

func BroadcastColorPlacements() {
	for placement := range broadcast {
		for client := range clients {
			select {
			case client.Send <- placement:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
		fmt.Println("Broadcasted color placement")
	}
}
