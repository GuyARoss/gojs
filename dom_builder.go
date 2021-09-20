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

					const handleIncomingRequest = ({ name, data }) => {
						switch (name) {
							case "register_event": {
								const { id, eventName } = data

								// adds an event listener to the a dom element
								document.getElementById(id).addEventListener(eventName, () => {
									// @@todo: create a copy of the dom & pass it into the emit event.
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
							const incomingJSON = JSON.parse(incoming.data)
							handleIncomingRequest(incomingJSON)
						}
					})
				</script>
			</head>
			<body>
			</body>
		</html>
	`, b.SocketURI)
}
