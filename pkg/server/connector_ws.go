package server

import (
	"encoding/json"
	"net/http"

	"github.com/dobin/antnium/pkg/campaign"
	"github.com/dobin/antnium/pkg/model"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ConnectorWs struct {
	middleware *Middleware
	clients    map[string]*websocket.Conn // ComputerId:WebsocketConnection
	wsUpgrader websocket.Upgrader
	coder      model.Coder
	campaign   *campaign.Campaign
}

func MakeConnectorWs(campaign *campaign.Campaign, middleware *Middleware) ConnectorWs {
	a := ConnectorWs{
		middleware: middleware,
		clients:    make(map[string]*websocket.Conn),
		wsUpgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		coder:      model.MakeCoder(campaign),
		campaign:   campaign,
	}
	return a
}

func (cw ConnectorWs) Shutdown() {
	for _, conn := range cw.clients {
		if conn != nil {
			conn.Close()
		}
	}
}

// wsHandlerClient is the entry point for new client initiated websocket connections
func (a *ConnectorWs) wsHandlerClient(w http.ResponseWriter, r *http.Request) {
	ws, err := a.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("ClientWebsocket: %s", err.Error())
		return
	}

	// WebSocket Authentication
	var authToken model.ClientWebSocketAuth
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Error("ClientWebsocket read error")
		return
	}
	err = json.Unmarshal(message, &authToken)
	if err != nil {
		log.Errorf("ClientWebsocket: could not decode auth: %v", message)
		return
	}
	if authToken.Key != "antnium" {
		log.Warn("ClientWebsocket: incorrect key: " + authToken.Key)
		return
	}
	// register client as auth succeeded
	a.clients[authToken.ComputerId] = ws

	a.handleWs(authToken.ComputerId, ws)
}

func (a *ConnectorWs) handleWs(computerId string, ws *websocket.Conn) {
	if ws == nil {
		log.Error("handleWs with invalid websocket connection")
		return
	}

	// Thread which reads from the client connection
	// Lifetime: Websocket connection
	go func() {
		for {
			_, packetData, err := ws.ReadMessage()
			if err != nil {
				ws.Close()
				a.clients[computerId] = nil
				break
			}
			packet, err := a.coder.DecodeData(packetData)
			if err != nil {
				log.Infof("registerWs error: %s", err.Error())
				continue
			}

			a.middleware.ClientSendPacket(packet, ws.RemoteAddr().String(), "ws")
		}
	}()

	// send all packets which havent yet been answered
	// make sure its a copy, and only iterate once.
	// If server is not available (WS disconnected), the packet response is lost.

	// make it a thread, so we return and all the stuff works
	//go func() {
	packets := make([]model.Packet, 0)
	for {
		packet, ok := a.middleware.ClientGetPacket(computerId, ws.RemoteAddr().String(), "ws")
		if !ok {
			break
		}
		packets = append(packets, packet)
	}
	for _, packet := range packets {
		ok := a.TryViaWebsocket(&packet)
		if !ok {
			log.Warn("Sending of initial packets via websocket failed")
		}
	}
	//}()

}

func (a *ConnectorWs) TryViaWebsocket(packet *model.Packet) bool {
	clientConn, ok := a.clients[packet.ComputerId]
	if !ok {
		// All ok, not connected to ws
		return false
	}
	if clientConn == nil {
		log.Warn("WS Client connection nil")
		return false
	}

	// Encode the packet and send it
	jsonData, err := a.coder.EncodeData(*packet)
	if err != nil {
		return false
	}

	err = clientConn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Infof("Websocket for host %s closed when trying to write: %s", packet.ComputerId, err.Error())
		return false
	}

	log.Debugf("Sent packet %s to client %s via WS", packet.PacketId, packet.ComputerId)

	return true
}
