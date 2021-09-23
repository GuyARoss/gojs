package gojs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/pkgz/websocket"
)

type DevServer struct {
	wsServer *websocket.Server
	started  bool

	operations *[]func(c *websocket.Conn)
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

	if s.started {
		s.wsServer.Emit("set_element", f)
		return
	}

	ops := append(*s.operations, func(c *websocket.Conn) {
		s.wsServer.Emit("set_element", f)
	})

	s.operations = &ops
}
func (s *DevServer) Setup() {
	s.wsServer.OnConnect(func(c *websocket.Conn) {
		fmt.Println(len(*s.operations))
		for _, op := range *s.operations {
			op(c)
		}
	})
}

func (s *DevServer) RenderDOM(body string) {
	// @@todo(guy): can this be used after the server is started?
	ops := append(*s.operations, func(c *websocket.Conn) {
		c.Emit("render_dom", body)
	})

	s.operations = &ops
}

func (s *DevServer) RegisterEventBridge() *UIUpdate {
	elChan := make(chan EventListenerEvent)

	s.wsServer.OnMessage(func(c *websocket.Conn, h ws.Header, b []byte) {
		fmt.Println("received", string(b))
	})

	return &UIUpdate{
		EventListenerSignal: elChan,
	}
}

func NewDevServer(serverPort int) UIClient {
	wsServer := websocket.Start(context.Background())
	ops := make([]func(c *websocket.Conn), 0)

	http.HandleFunc("/ws", wsServer.Handler)

	go func() {
		http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			dom := &DOMBuilder{
				SocketURI: fmt.Sprintf("ws://localhost:%d/ws", serverPort),
			}
			rw.Write([]byte(dom.Build()))
			rw.WriteHeader(http.StatusOK)
		})

		http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
	}()

	return &DevServer{
		wsServer:   wsServer,
		operations: &ops,
	}
}
