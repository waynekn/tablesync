import { X } from "lucide-react";
import { Trash2 } from "lucide-react";
import { Button } from "./ui/button";

type ColHeadersListProps = {
  columns: string[];
  isSidebarOpen: boolean;
  toggleSidebar: () => void;
  rmHeader: (index: number) => void;
};

export default function ColTitlesList({
  columns,
  toggleSidebar,
  isSidebarOpen,
  rmHeader,
}: ColHeadersListProps) {
  return (
    <aside
      className={`md:w-1/4 w-full bg-white border-r-2 border-gray-300 shadow-lg p-4 md:static fixed inset-y-0 left-0 z-10 transform transition-transform duration-300 ${
        isSidebarOpen ? "translate-x-0" : "-translate-x-full"
      } md:translate-x-0`}
      id="sidebar"
    >
      <div className="relative">
        <Button
          className="absolute top-2 right-2  md:hidden"
          onClick={toggleSidebar}
        >
          <X size={20} />
        </Button>
      </div>

      <p className="block text-gray-700 font-medium text-sm mt-6">
        Column Titles
      </p>

      <ol>
        {columns.map((column, index) => (
          <li key={index} className="flex list-decimal">
            <p key={index} className="grow">
              {column}
            </p>
            <Button onClick={() => rmHeader(index)} variant="ghost">
              <Trash2 />
            </Button>
          </li>
        ))}
      </ol>
    </aside>
  );
}
