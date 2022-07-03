import { useEffect, useMemo, useRef, useState } from "react";
import {
  useWindowSize,
  useKeyPress,
  useKeyPressEvent,
  useRafLoop,
} from "react-use";
import "./App.css";

function App() {
  const { width, height } = useWindowSize();
  const position = useRef({ x: width / 2, y: height / 2 });
  const keys = useRef({
    ArrowLeft: false,
    ArrowRight: false,
    ArrowUp: false,
    ArrowDown: false,
  });

  const setKeyDown = (key: keyof typeof keys["current"]) => () => {
    keys.current[key] = true;
  };

  const setKeyUp = (key: keyof typeof keys["current"]) => () => {
    keys.current[key] = false;
  };

  for (const key in keys.current) {
    const k = key as keyof typeof keys["current"];
    useKeyPressEvent(k, setKeyDown(k), setKeyUp(k));
  }

  const canvasRef = useRef<HTMLCanvasElement>(null);
  const context = useRef<CanvasRenderingContext2D | null>(null);

  const gameLoop = (time: number) => {
    const moveSpeed = 5;
    if (keys.current.ArrowLeft) {
      position.current.x -= 5;
    }
    if (keys.current.ArrowRight) {
      position.current.x += 5;
    }
    if (keys.current.ArrowUp) {
      position.current.y -= 5;
    }
    if (keys.current.ArrowDown) {
      position.current.y += 5;
    }
  };

  useRafLoop(gameLoop);

  const draw = (ctx: CanvasRenderingContext2D) => {
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
    ctx.fillStyle = "#000000";
    ctx.beginPath();
    ctx.arc(position.current.x, position.current.y, 20, 0, 2 * Math.PI);
    ctx.fill();
  };

  useRafLoop((time) => {
    if (!context.current && canvasRef.current) {
      context.current = canvasRef.current.getContext("2d");
    }

    if (context.current) {
      draw(context.current);
    }
  });

  return <canvas ref={canvasRef} width={width} height={height}></canvas>;
}

export default App;
