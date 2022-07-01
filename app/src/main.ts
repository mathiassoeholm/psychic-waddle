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
  const message = new TextEncoder().encode("Hello other players!");
  websocket.send(Uint8Array.of(MessageId.SendChatMessage, ...message));
});

websocket.addEventListener("error", (error) => {
  console.error(error);
});

websocket.addEventListener("message", (message) => {
  console.log("Got message data", message.data);
  if (message.data[0] === MessageId.ReceiveChatMessage) {
  }
});
