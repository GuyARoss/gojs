package gojs

import "fmt"

type DOMBuilder struct {
	SocketURI string
}

func (b *DOMBuilder) Build() string {
	return fmt.Sprintf(`
		<html>			
			<head>
				<script>
					const socket = new WebSocket("%s")

					const staticDoc = () => {
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

					const handleIncomingRequest = ({ name, data }) => {
						switch (name) {
							case "register_event": {
								const { id, eventName } = data

								// adds an event listener to the a dom element
								document.getElementById(id).addEventListener(eventName, () => {
									const staticDoc = staticDoc()
									console.log(eventName, id)
								})
							}
							case "set_element": {
								const { elementID, content } = data

								// sets the inner html of an element
								const element = document.getElementById(elementID)
								element.innerHTML = content
							}
							case "render_dom": {
								document.body.innerHTML = data
							}
						}
					}

					socket.addEventListener("message", (e) => {
						const incoming = JSON.parse(e.data)
			
						if (!!incoming.data) {
							handleIncomingRequest(incoming)
						}
					})
				</script>
			</head>
			<body>
			</body>
		</html>
	`, b.SocketURI)
}
