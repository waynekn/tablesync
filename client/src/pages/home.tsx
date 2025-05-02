import { useAuth } from "react-oidc-context";
import { Button } from "@/components/ui/button";

function HomePage() {
  const auth = useAuth();
  document.title = "Home";

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4 text-center">
      <div className="max-w-2xl">
        <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-6">
          Real-Time Collaboration on Spreadsheets
        </h1>
        <p className="text-lg md:text-xl text-gray-700 mb-4">
          TableSync lets your team edit, organize, and collect data together for
          collaborative workflows.
        </p>
        <p className="text-md md:text-lg text-gray-600 mb-8">
          Boost productivity with live collaboration and built-in data
          collection features.
        </p>
        <Button
          onClick={() => auth.signinRedirect()}
          variant="dark"
          disabled={auth.isLoading}
        >
          Get Started
        </Button>
      </div>
    </div>
  );
}

export default HomePage;
