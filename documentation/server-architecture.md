```mermaid
flowchart TB
  websockets[Gorilla Websocket]

  websockets -- r_chan --> ws_server
  ws_server -- w_chan --> websockets

  ws_server[Websocket Server]

  app[App]

  ws_server -- chan --> app
  app -- chan --> ws_server
```

### Messages in channels:

Websocket Server to App:

- Player connected
- Player sent message
- Player disconnected

App to Websocket Server:

- Emit message to player
- Disconnect player

r_chan:

- Player sent bytes

w_chan:

- Send bytes to player
