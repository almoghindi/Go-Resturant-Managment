import { useState, useEffect, useCallback } from "react";

interface User {
  user_id: string;
  firstName: string;
  email: string;
  token: string;
  refresh_token: string;
}

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
  }, []);

  const login = useCallback(
    async (
      user_id: string,
      email: string,
      firstName: string,
      token: string,
      refresh_token: string
    ) => {
      try {
        const userData = {
          user_id,
          email,
          firstName,
          token,
          refresh_token,
        };
        setUser(userData);
        localStorage.setItem("user", JSON.stringify(userData));
      } catch (error) {
        console.error("Login failed", error);
        throw new Error("Login failed");
      }
    },
    []
  );

  const logout = useCallback(() => {
    setUser(null);
    localStorage.removeItem("user");
  }, []);

  return { user, login, logout };
};
