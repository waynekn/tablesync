// EditSheet.tsx
import { useRef, useState, useCallback } from "react";
import { useParams } from "react-router";
import { registerAllModules } from "handsontable/registry";
import { HotTable, HotTableRef } from "@handsontable/react-wrapper";

import "handsontable/styles/handsontable.min.css";
import "handsontable/styles/ht-theme-main.min.css";

import { useSheetWebSocket } from "@/hooks/useSheetWebSocket";
import getContextMenuSettings from "@/lib/contextMenuSettings";

registerAllModules();

export default function EditSheet() {
  const { sheetID } = useParams();
  const hotTableRef = useRef<HotTableRef>(null);
  const skipNextChange = useRef(false);

  const [sheetData, setSheetData] = useState<string[][]>([[]]);
  const [tableColumns, setTableColumns] = useState<string[]>([]);

  const handleInitData = useCallback((data: string[][], cols: string[]) => {
    setSheetData(data);
    setTableColumns(cols);
  }, []);

  const handleRemoteEdit = useCallback(
    (row: number, col: number, data: string) => {
      skipNextChange.current = true;
      hotTableRef.current?.hotInstance?.setDataAtCell(row, col, data);
    },
    []
  );

  const { sendEdit } = useSheetWebSocket(
    sheetID,
    handleInitData,
    handleRemoteEdit
  );

  return (
    <>
      <p>Live edit</p>
      <div className="ht-theme-main">
        <HotTable
          ref={hotTableRef}
          data={sheetData}
          rowHeaders
          colHeaders={tableColumns.length > 0 ? tableColumns : true}
          height="auto"
          autoWrapRow
          autoWrapCol
          contextMenu={getContextMenuSettings(hotTableRef)}
          enterMoves={{ row: 0, col: 1 }}
          afterChange={(changes, source) => {
            if (skipNextChange.current) {
              skipNextChange.current = false;
              return;
            }
            if (source === "edit" && changes) {
              changes.forEach(([r, c, , newValue]) => {
                sendEdit(r, Number(c), newValue as string);
              });
            }
          }}
          licenseKey="non-commercial-and-evaluation"
        />
      </div>
    </>
  );
}
