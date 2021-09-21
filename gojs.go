package gojs

import (
	"io/ioutil"
)

type StaticDocument struct {
	// ContentMap where the key is the ID of the element
	// and the value is the content within the Element
	ContentMap map[string]string
}

type EventListenerEvent struct {
	ElementID string
	EventName string
	Document  *StaticDocument
}

type UIUpdate struct {
	EventListenerSignal chan EventListenerEvent
}

type UIClient interface {
	RegisterEvent(id string, eventName string)
	SetElement(elementID string, data string)
	RenderDOM(body string)
	RegisterEventBridge() *UIUpdate
	Setup()
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

func (ui *UI) Show() {
	ui.Client.Setup()
	pathContent, err := ioutil.ReadFile(ui.Config.HTMLDocPath)
	if err != nil {
		panic("invalid doc path")
	}

	// @@todo(guy): strip everything besides body tags
	ui.Client.RenderDOM(string(pathContent))

	bridge := ui.Client.RegisterEventBridge()
	finished := make(chan bool)

	go func() {
		for {
			select {
			case s := <-bridge.EventListenerSignal:
				go ui.events[s.ElementID][s.EventName](s.Document)
			}
		}
		// @@ todo(guy): this goes on forever, need to close it..
	}()

	<-finished
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
