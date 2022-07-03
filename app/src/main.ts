import "./style.css";
import { MessageId } from "./protocol";

const app = document.querySelector<HTMLDivElement>("#app")!;

app.innerHTML = `
  <h1>Hello Vite!</h1>
  <a href="https://vitejs.dev/guide/features.html" target="_blank">Documentation</a>
`;

const websocket = new WebSocket("ws://localhost:4000");
websocket.binaryType = "arraybuffer";
websocket.addEventListener("open", () => {
  console.log("Connected to websocket");
  const message = new TextEncoder().encode("Hello other players! 🦙");
  websocket.send(Uint8Array.of(MessageId.SendChatMessage, ...message));
});

websocket.addEventListener("error", (error) => {
  console.error(error);
});

websocket.addEventListener("message", (message) => {
  console.log("Got message with id", new Uint8Array(message.data)[0]);
  const data = new Uint8Array(message.data);
  if (data[0] === MessageId.ReceiveChatMessage) {
    const playerId = new DataView(data.slice(1, 5).buffer).getUint32(0);
    const message = new TextDecoder().decode(data.slice(5));
    console.log({ playerId, message });
  }
});
