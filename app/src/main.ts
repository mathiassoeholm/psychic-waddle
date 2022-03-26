import "./style.css";

const app = document.querySelector<HTMLDivElement>("#app")!;

app.innerHTML = `
  <h1>Hello Vite!</h1>
  <a href="https://vitejs.dev/guide/features.html" target="_blank">Documentation</a>
`;

const websocket = new WebSocket("ws://localhost:4000");
websocket.addEventListener("open", () => {
  console.log("Connected to websocket");
  websocket.send("Hello Server!");
});

websocket.addEventListener("error", (error) => {
  console.error(error);
});

websocket.addEventListener("message", (message) => {
  console.log(message);
});
