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

							}
							case "set_element": {
								
							}
							case "render_dom": {
								
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
