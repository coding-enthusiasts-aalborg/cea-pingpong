package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var broadcast = make(chan message)
var games = make(map[string]map[string]player)

const playersPerGame = 2

type message struct {
	GameID   string //`json:"gameid"`
	PlayerID string //`json:"playerid"`
	Y        float32
}

type player struct {
	conn *websocket.Conn
}

// Configure upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r := mux.NewRouter()

	r.HandleFunc("/ws", handleConnections).Queries("gameid", "{gameid}", "playerid", "{playerid}")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	go handleMessages()

	log.Printf("Running server on http://localhost:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	vars := mux.Vars(r)
	gameID := vars["gameid"]
	playerID := vars["playerid"]
	newPlayer := player{ws}

	fmt.Println("GameID:", gameID, "PlayerID", playerID)

	// Make sure we close the connection when the function returns
	defer ws.Close()

	game, ok := games[gameID]
	if !ok {
		games[gameID] = make(map[string]player)
		game = games[gameID]
	}
	if len(game) >= playersPerGame {
		return
	}

	games[gameID][playerID] = newPlayer

	if len(game) == playersPerGame {
		go sendStartGameMessages(games[gameID])
	}

	for {
		var msg message
		// Read the new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		fmt.Println(msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(games[gameID], playerID)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func sendStartGameMessages(players map[string]player) {
	positions := []string{"left", "right"}
	i := 0
	for _, p := range players {
		p.conn.WriteJSON("{\"position\": \"" + positions[i] + "\"}")
		i++
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		fmt.Println("GameID: ", msg.GameID, "PlayerID: ", msg.PlayerID, "Y: ", msg.Y)

		// Send it to the all the players in the specified game
		for playerID, player := range games[msg.GameID] {

			// Don't send msg back to sender
			if playerID == msg.PlayerID {
				continue
			}

			err := player.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				delete(games[msg.GameID], playerID)
			}
		}
	}
}
