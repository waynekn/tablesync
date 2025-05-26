import { useEffect, useState } from "react";
import { AxiosError } from "axios";
import { Link } from "react-router";
import { User } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { useAuth, withAuthenticationRequired } from "react-oidc-context";

import { useApi } from "@/hooks/api";

import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import UserInfoCard from "@/components/cards/userInfo";
import SheetList from "@/components/sheetList";
import ListSkeleton from "@/skeletons/listSkeleton";

import { GenericApiError } from "@/types/error";
import { Spreadsheet } from "@/types/spreadsheet";

function ProfilePage() {
  const auth = useAuth();
  const api = useApi();
  const [sheets, setSheets] = useState<Spreadsheet[]>([]);
  const [errMsg, setErrMsg] = useState<string>("");

  const { data, error, isSuccess, isPending, isError } = useQuery({
    queryKey: [`${auth.user?.profile.sub}`],
    queryFn: () => api.get<Spreadsheet[]>("spreadsheets/"),
    retry: (failureCount, error) => {
      const apiError = error as AxiosError;
      if (apiError.response?.status === 500) return false; // No retry on server errors
      return failureCount < 3;
    },
    staleTime: 1000 * 60 * 30, // 30 minutes
  });

  useEffect(() => {
    if (isSuccess) {
      setSheets(data.data);
      return;
    }

    if (isError) {
      const err = error as AxiosError<GenericApiError>;
      setErrMsg(err.response?.data.error || "An unexpected error occurred");
    }
  }, [data, error, isError, isSuccess]);

  return (
    <div className="h-screen w-screen flex flex-col md:flex-row bg-gray-100">
      <main className="flex-1 p-4 md:p-8 sm:text-sm md:text-base">
        <Popover>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              size="icon"
              className="rounded-full border-gray-300 hover:bg-gray-100 hover:border-gray-400 transition-colors"
            >
              <User className="h-5 w-5 text-gray-800" />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-80 border-gray-200 bg-white text-gray-900 shadow-md rounded-lg p-0">
            <UserInfoCard />
          </PopoverContent>
        </Popover>

        <section className="w-full max-w-3xl mx-auto border border-gray-200 rounded-xl bg-white shadow-md p-10 space-y-8">
          <div>
            <Link
              to="../spreadsheet/create/"
              className="inline-block text-blue-600 font-semibold underline hover:text-blue-800 transition-colors duration-200"
            >
              Create a New Spreadsheet
            </Link>
          </div>

          {isPending && <ListSkeleton />}

          {isError && (
            <div
              role="alert"
              className="bg-red-600 text-white rounded-lg py-4 px-2"
            >
              {errMsg || "An unexpected error occurred. Please try again."}
            </div>
          )}

          {isSuccess && (
            <div>
              {sheets.length > 0 ? (
                <SheetList sheets={sheets} />
              ) : (
                <p className="text-gray-500">No spreadsheets found.</p>
              )}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

export default withAuthenticationRequired(ProfilePage);
