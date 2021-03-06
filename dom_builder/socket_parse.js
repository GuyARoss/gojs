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
