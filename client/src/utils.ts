import { SheetEdit } from "./types/webSocket";

export function is2DArray(value: unknown): value is [][] {
  return Array.isArray(value) && value.every((item) => Array.isArray(item));
}

export function isSheetEdit(value: unknown): value is SheetEdit {
  return (
    value !== null &&
    typeof value === "object" &&
    typeof (value as SheetEdit).row === "number" &&
    typeof (value as SheetEdit).col === "number" &&
    typeof (value as SheetEdit).data === "string"
  );
}
