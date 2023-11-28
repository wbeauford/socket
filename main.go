package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os/exec"
)
// two d's for a double dose of pimpin'
var upgrayedd = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	upgrayedd.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade connection
	ws, err := upgrayedd.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Connected")


	//send message back
	err = ws.WriteMessage(1, []byte("What up"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

// reader function to listen to new messages sent
func reader(conn *websocket.Conn) {
	for {
		messageType , p , err := conn.ReadMessage()
		if err != nil {
			// error reading message
			log.Println(err)
			return
		}
		// Print out the message
		log.Println(string(p))
		out := runCommand(string(p))
		log.Println(out)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func runCommand(command string) string {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleWS)
	server :=&http.Server{
		Handler: r,
		Addr: ":8080",
	}

	log.Fatal(server.ListenAndServe())
}




