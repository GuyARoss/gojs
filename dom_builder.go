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

					const handleIncomingRequest = ({ name, data }) => {
						switch (name) {
							case "register_event": {
								const { id, eventName } = JSON.parse(data)

								document.getElementById(id).addEventListener(eventName, () => {
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
								break;
							}
							case "set_element": {
								const { elementID, content } = JSON.parse(data)

								// sets the inner html of an element
								const element = document.getElementById(elementID)
								element.innerHTML = content
								break;

							}
							case "render_dom": {
								document.body.innerHTML = data
								break;

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
