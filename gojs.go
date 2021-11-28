package gojs

import (
	"errors"
	"io/ioutil"
	"time"
)

type StaticDocument struct {
	// ContentMap where the key is the ID of the element
	// and the value is the content within the Element
	ContentMap map[string]string `json:"contentMap"`
}

type EventListenerEvent struct {
	ElementID string          `json:"elementID"`
	EventName string          `json:"eventName"`
	Document  *StaticDocument `json:"document"`
}

var ErrHostConnectionTerminated = errors.New("host connection terminated")

type UIUpdate struct {
	EventListenerSignal chan EventListenerEvent
	EventErrorSignal    chan error
}

type UIClient interface {
	RegisterEvent(id string, eventName string)
	SetElement(elementID string, data string)
	RenderDOM(body string)
	RegisterEventBridge() *UIUpdate
	Setup()
	IsActiveConnection() bool
}

type StaticDOMElement interface {
	GetContent() string
}

type EventHandler func(doc *StaticDocument)

type UIConfig struct {
	HTMLDocPath string
}

type UI struct {
	// events, where the primary key is the name of the
	// element & the secondary key is the name of the event name
	events map[string]map[string]EventHandler

	Client UIClient
	Config *UIConfig
}

func New(client UIClient, config *UIConfig) *UI {
	return &UI{
		events: make(map[string]map[string]EventHandler),
		Client: client,
		Config: config,
	}
}

func (ui *UI) Show() {
	ui.Client.Setup()
	pathContent, err := ioutil.ReadFile(ui.Config.HTMLDocPath)

	// @@todo(guy): strip everything besides body tags
	ui.Client.RenderDOM(string(pathContent))

	if err != nil {
		panic("invalid doc path")
	}

	bridge := ui.Client.RegisterEventBridge()
	terminationErrChan := make(chan error)

	go func() {
		for {
			select {
			case s := <-bridge.EventListenerSignal:
				go ui.events[s.ElementID][s.EventName](s.Document)
			case err := <-bridge.EventErrorSignal:
				terminationErrChan <- err
			}
		}
	}()

	// ensure that the connection exists before we do stuff.
	connectionEstablished := make(chan bool)
	go func() {
		for {
			time.Sleep(time.Millisecond * 50)
			if ui.Client.IsActiveConnection() {
				connectionEstablished <- true
				break
			}
		}
	}()

	// we wait for the connection to be established before attempting to render anything.
	<-connectionEstablished

	terminationErr := <-terminationErrChan
	if errors.Is(terminationErr, ErrHostConnectionTerminated) {
		// @@todo(guy): reset?
	}
}

type StaticDOMElementInstance struct {
	ElementID string

	doc *StaticDocument
}

func (s *StaticDOMElementInstance) GetContent() string {
	return s.doc.ContentMap[s.ElementID]
}

func (d *StaticDocument) Element(elementID string) StaticDOMElement {
	return &StaticDOMElementInstance{
		ElementID: elementID,
		doc:       d,
	}
}

type LiveDOMElement interface {
	SetEventListener(key string, fn EventHandler)
	SetContent(data string)
}

type LiveDOMElementInstance struct {
	ui        *UI
	ElementID string
}

func (r *LiveDOMElementInstance) SetEventListener(key string, fn EventHandler) {
	if r.ui.events[r.ElementID] == nil {
		r.ui.events[r.ElementID] = make(map[string]EventHandler)
	}

	r.ui.events[r.ElementID][key] = fn

	// this should ensure that the event exists on the clients element i.e create it if not exists.
	go r.ui.Client.RegisterEvent(r.ElementID, key)
}

func (r *LiveDOMElementInstance) SetContent(data string) {
	go r.ui.Client.SetElement(r.ElementID, data)
}

func (ui *UI) Element(key string) LiveDOMElement {
	return &LiveDOMElementInstance{
		ui:        ui,
		ElementID: key,
	}
}
