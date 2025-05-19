import { User } from "lucide-react";
import { withAuthenticationRequired } from "react-oidc-context";

import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import UserInfoCard from "@/components/cards/userInfo";
import { Link } from "react-router";

function ProfilePage() {
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

        <section className="w-full max-w-3xl mx-auto border-2 border-gray-300 rounded-lg bg-white shadow-lg p-8 space-y-6">
          <div>
            <Link
              to="../spreadsheet/create/"
              className="text-black underline hover:text-gray-700"
            >
              Create spreadsheet
            </Link>
          </div>
          {/* TODO */}
        </section>
      </main>
    </div>
  );
}

export default withAuthenticationRequired(ProfilePage);
