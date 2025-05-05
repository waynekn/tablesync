import { useAuth } from "react-oidc-context";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

/**
 * While this card displays user information, it is not wrapped in `react-oidc-context's`
 * `withAuthenticationRequired` because it is currently only used on
 * the profile page, which already has authentication enforced.
 */
export default function UserInfoCard() {
  const auth = useAuth();

  return (
    <Card className="border-0 bg-transparent">
      <CardHeader className="border-b border-gray-200">
        <CardTitle className="text-lg font-medium text-gray-900">
          Profile Info
        </CardTitle>
      </CardHeader>
      <CardContent className="text-sm">
        <p className="mb-2">
          <span className="font-semibold text-gray-700">Username:</span>
          {auth.user?.profile.given_name || "N/A"}
        </p>
        <p>
          <span className="font-semibold text-gray-700">Email:</span>
          {auth.user?.profile.email || "N/A"}
        </p>
      </CardContent>
    </Card>
  );
}
