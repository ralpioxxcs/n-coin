package p2p

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ralpioxxcs/n-coin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// port :3000 will upgrade the request from :4000
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// openPort is used by other peers to request each other peers
		return openPort != "" && ip != ""
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)

	peer := initPeer(conn, ip, openPort)
	time.Sleep(20 * time.Second)
	peer.inbox <- []byte("Hi from 3000!")
}

/*
	+----------+	  	+----------+
	|  :4000   | ---> |	 :3000   |
	+----------+	  	+----------+
	- add peer
	- make websocket connection to port 3000
*/
func AddPeer(address, port, openPort string) {
	// Port :4000 is requesting an upgrade from the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)

	peer := initPeer(conn, address, port)
	time.Sleep(10 * time.Second)
	peer.inbox <- []byte("Hello from port 4000!")
}
