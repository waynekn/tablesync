// useSheetWebSocket.ts
import { useEffect, useRef } from "react";
import { SheetEdit } from "@/types/webSocket";
import { is2DArray, isSheetEdit } from "@/utils";

export function useSheetWebSocket(
  sheetID: string | undefined,
  onInitData: (data: string[][], columns: string[]) => void,
  onEdit: (row: number, col: number, data: string) => void
) {
  const socket = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!sheetID) return;

    socket.current = new WebSocket(
      `ws://localhost:8000/ws/sheet/${sheetID}/edit/`
    );

    socket.current.onmessage = (e) => {
      const msg = JSON.parse(e.data);
      if (is2DArray(msg)) {
        const cols = msg.shift() as string[];
        const payload =
          msg.length === 0 ? [new Array(cols.length).fill("")] : msg;
        onInitData(payload, cols);
      } else if (isSheetEdit(msg)) {
        onEdit(msg.row, msg.col, msg.data);
      }
    };

    return () => {
      socket.current?.close(
        1000,
        "Your connection has been closed. Please refresh to connect"
      );
      socket.current = null;
    };
  }, [sheetID, onInitData, onEdit]);

  function sendEdit(row: number, col: number, data: string) {
    if (socket.current && socket.current.readyState === WebSocket.OPEN) {
      const edit: SheetEdit = { row, col, data };
      socket.current.send(JSON.stringify(edit));
    }
  }

  return { sendEdit };
}
