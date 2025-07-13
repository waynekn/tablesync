import { ContextMenu, Settings } from "handsontable/plugins/contextMenu";
import { HotTableRef } from "@handsontable/react-wrapper";

export default function getContextMenuSettings(
  hotTableRef: React.RefObject<HotTableRef | null>
): Settings {
  if (!hotTableRef.current) {
    return {};
  }

  return {
    items: {
      row_below: {
        callback: function () {
          const hotInstance = hotTableRef.current?.hotInstance;
          if (hotInstance) {
            const totalRows = hotInstance.countRows();
            hotInstance.alter("insert_row_below", totalRows);
          }
        },
        disabled: function () {
          const hotInstance = hotTableRef.current?.hotInstance;
          if (!hotInstance) return false;
          const selected = hotInstance.getSelectedLast();
          const totalRows = hotInstance.countRows();
          if (!selected) return true;
          const selectedRow = selected[0];
          const lastRowIndex = totalRows - 1;
          return selectedRow !== lastRowIndex;
        },
      },
      sep1: ContextMenu.SEPARATOR,
      alignment: {},
      sep2: ContextMenu.SEPARATOR,
      copy: {},
      cut: {},
    },
  };
}
