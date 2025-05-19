import axios from "axios";
import { useAuth } from "react-oidc-context";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// This hook attaches the access token to the request
export function useApi() {
  const auth = useAuth();

  api.interceptors.request.use(
    (config) => {
      const token = auth.user?.access_token;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  return api;
}
