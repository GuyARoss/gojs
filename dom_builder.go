// Code generated by gojs/dom-builder. DO NOT EDIT.
// source dom_builder.go

package gojs

import "fmt"

type DOMBuilder struct {
	SocketURI    string
	HTMLHeadTags string
}

func (b *DOMBuilder) Build() string {
	return fmt.Sprintf(`
            <html>			
                <head>
                    %s

                    <script>
                        const socket = new WebSocket("%s")

                        function createStaticDoc() {
    const uiElements = document.getElementsByClassName("ui")
    const uiDoc = {}

    for (let i = 0; i < uiElements.length; i++) {
        const uiElement = uiElements.item(i)
        const uiElementID = uiElement.id

        if (!!uiElementID) {
            uiDoc[uiElementID] = uiElement.innerText
        }
    }

    return uiDoc
}

function handleIncomingRequest (socket, { name, data }) {
    switch (name) {
        case "register_event": {
            const { id, eventName } = JSON.parse(data)

            const f = document.getElementById(id)
            if (f !== null) {
                f.addEventListener(eventName, () => {
                    const staticDoc = createStaticDoc()

                    socket.send(JSON.stringify({
                        name: "event",
                        data: {
                            document: {
                                contentMap: staticDoc,
                            },
                            eventName: eventName,
                            elementID: id,
                        },
                    }))
                })
            }

            break;
        }
        case "set_element": {
            const { elementID, content } = JSON.parse(data)

            // sets the inner html of an element
            const element = document.getElementById(elementID)
            if (element !== null) {
                element.innerHTML = content
            }

            break;
        }
        case "render_dom": {
            document.body.innerHTML = data
            break;

        }
    }
}

function initializeSocket (socket) {
    socket.addEventListener("message", (e) => {
        const incoming = JSON.parse(e.data)
    
        if (!!incoming.data) {
            handleIncomingRequest(socket, incoming)
        }
    })
}


                        initializeSocket(socket)
                    </script>
                </head>
                <body>
                </body>
            </html>
        `, b.HTMLHeadTags, b.SocketURI)
}
