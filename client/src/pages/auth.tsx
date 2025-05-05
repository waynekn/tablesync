import { useAuth } from "react-oidc-context";
import { useNavigate } from "react-router";

function AuthPage() {
  const auth = useAuth();
  const navigate = useNavigate();

  const signOutRedirect = () => {
    const clientId = import.meta.env.VITE_CLIENT_ID;
    const logoutUri = "<logout uri>";
    const cognitoDomain = import.meta.env.VITE_COGNITO_DOMAIN;
    window.location.href = `${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(
      logoutUri
    )}`;
  };

  if (auth.isAuthenticated) {
    navigate(`../profile/${auth.user?.profile.given_name}`);
  }

  if (auth.isLoading) {
    return <div>Loading...</div>;
  }

  if (auth.error) {
    return <div>Encountering error... {auth.error.message}</div>;
  }

  if (auth.isAuthenticated) {
    return (
      <div>
        <p>Sucessfully authenticated</p>
      </div>
    );
  }

  return (
    <div>
      <button onClick={() => auth.signinRedirect()}>Sign in</button>
      <button onClick={() => signOutRedirect()}>Sign out</button>
    </div>
  );
}

export default AuthPage;
