package gojs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkgz/websocket"
)

type DevServer struct {
	wsServer *websocket.Server
}

func (s *DevServer) RegisterEvent(id string, eventName string) {
	m := make(map[string]string)
	m["id"] = id
	m["eventName"] = eventName

	f, _ := json.Marshal(m)
	s.wsServer.Emit("register_event", f)
}

func (s *DevServer) SetElement(elementID string, data string) {
	m := make(map[string]string)
	m["elementID"] = elementID
	m["content"] = data

	f, _ := json.Marshal(m)

	s.wsServer.Emit("render_dom", f)
}

func (s *DevServer) RenderDOM(body string) {
	s.wsServer.Emit("render_dom", []byte(body))
}

func (s *DevServer) RegisterEventBridge() *UIUpdate {
	return &UIUpdate{}
}

func NewDevServer() UIClient {
	wsServer := websocket.Start(context.Background())

	go func() {
		http.HandleFunc("/ws", wsServer.Handler)
		http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			dom := &DOMBuilder{}
			rw.Write([]byte(dom.Build()))
			rw.WriteHeader(http.StatusOK)
		})

		http.ListenAndServe(fmt.Sprintf(":%d", 3000), nil)
	}()

	return &DevServer{
		wsServer: wsServer}
}
