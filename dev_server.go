package gojs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkgz/websocket"
)

type DevServer struct {
	wsServer *websocket.Server
}

func (s *DevServer) RegisterEvent(id string, eventName string) {

}

func (s *DevServer) SetElement(elementID string, data string) {

}

func (s *DevServer) RenderDOM(body string) {

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
		wsServer: wsServer,
	}
}
