import { Skeleton } from "@/components/ui/skeleton";

/**
 * A React component that renders a skeleton loading state for a list with `count` items.
 * Uses the `Skeleton` component from shadcn/ui to display placeholder list items.
 *
 * @param
 * - count - The number of list items skeletons to display
 * - height - The height of a list item
 * - bgColor - The background color of a list item
 */
export default function ListSkeleton({
  count = 10,
  height = "h-8",
  bgColor = "bg-gray-300",
}) {
  return (
    <div className="space-y-4 overflow-y-hidden">
      <ul className="space-y-3">
        {[...Array(count)].map((_, i) => (
          <li key={i}>
            <Skeleton className={`${height} w-full ${bgColor}`} />
          </li>
        ))}
      </ul>
    </div>
  );
}
