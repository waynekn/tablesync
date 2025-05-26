import { Spreadsheet } from "@/types/spreadsheet";

type SheetListProps = {
  sheets: Spreadsheet[];
};
/**
 * Displays an orderd list of spreadsheets
 */
export default function SheetList({ sheets }: SheetListProps) {
  return (
    <div>
      <h2 className="text-lg font-semibold text-gray-800 mb-4">
        Your Spreadsheets
      </h2>
      <ol className="list-decimal list-inside space-y-2 text-gray-700">
        {sheets.map((sheet) => (
          <li key={sheet.id} className="pl-2">
            {sheet.title}
          </li>
        ))}
      </ol>
    </div>
  );
}
